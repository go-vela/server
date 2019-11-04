// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-vela/compiler/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PostWebhook represents the API handler to capture
// a webhook from a source control provider and
// publish it to the configure queue.
func PostWebhook(c *gin.Context) {
	logrus.Info("Webhook received")

	// process the webhook from the source control provider
	r, b, err := source.FromContext(c).ProcessWebhook(c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to parse webhook: %w", err)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// check if build was parsed from webhook
	if b == nil {
		// typically, this should only happen on a webhook
		// "ping" which gets sent when the webhook is created
		c.JSON(http.StatusOK, "no build to process")
		return
	}

	// check if repo was parsed from webhook
	if r == nil {
		retErr := fmt.Errorf("unable to process webhook: failed to parse repo from webhook")
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// send API call to capture parsed repo from webhook
	r, err = database.FromContext(c).GetRepo(r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to process webhook: failed to get repo %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// check if the repo is active
	if !r.GetActive() {
		retErr := fmt.Errorf("unable to process webhook: %s is not an active repo", r.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// verify the build has a valid event and the repo allows that event type
	if (b.GetEvent() == constants.EventPush && !r.GetAllowPush()) ||
		(b.GetEvent() == constants.EventPull && !r.GetAllowPull()) ||
		(b.GetEvent() == constants.EventTag && !r.GetAllowTag()) ||
		(b.GetEvent() == constants.EventDeploy && !r.GetAllowDeploy()) {
		retErr := fmt.Errorf("unable to process webhook: %s does not have %s events enabled", r.GetFullName(), b.GetEvent())
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// check if the repo has a valid owner
	if r.GetUserID() == 0 {
		retErr := fmt.Errorf("unable to process webhook: %s has no valid owner", r.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// send API call to capture repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to process webhook: failed to get owner for %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// send API call to capture the last build for the repo
	lastBuild, err := database.FromContext(c).GetLastBuild(r)
	if err != nil {
		retErr := fmt.Errorf("unable to process webhook: failed to get last build for %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	// update fields in build object
	num := 1
	pending := constants.StatusPending
	time := time.Now().UTC().Unix()
	b.RepoID = r.ID
	b.Status = &pending
	b.Created = &time
	b.Number = &num
	b.Parent = b.Number

	bNumber := (lastBuild.GetNumber() + 1)
	bParent := lastBuild.GetNumber()
	if lastBuild != nil {
		b.Number = &bNumber
		b.Parent = &bParent
	}

	// variable to store changeset files
	files := []string{}
	// check if the build event is pull_request
	if b.GetEvent() == constants.EventPull {

		// parse out pull request number from base ref
		// TODO: clean this up
		s := strings.Split(b.GetRef(), "pull/")
		number, _ := strconv.Atoi(strings.Split(s[1], "/")[0])

		// send API call to capture list of files changed for the pull request
		files, err = source.FromContext(c).ListChangesPR(u, r, number)
		if err != nil {
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)
			return
		}
	} else { // all other event types use commit to get changes
		// send API call to capture list of files changed for the commit
		files, err = source.FromContext(c).ListChanges(u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)
			return
		}
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).Config(u, r.GetOrg(), r.GetName(), b.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("unable to process webhook: failed to get pipeline configuration for %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusNotFound, retErr)
		return
	}

	// parse and compile the pipeline configuration file
	pipe, err := compiler.FromContext(c).
		WithBuild(b).
		WithFiles(files).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, fmt.Errorf("Error compiling pipeline configuration for %s: %v", r.GetFullName(), err))
		return
	}

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), pipe, b, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, b)

	// send API call to set the status on the commit
	err = source.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromContext(c),
		pipe,
		b,
		r,
		u,
	)
}

// publishToQueue is a helper function that creates
// a build item and publishes it to the queue.
func publishToQueue(queue queue.Service, p *pipeline.Build, b *library.Build, r *library.Repo, u *library.User) {
	item := types.ToItem(p, b, r, u)

	logrus.Infof("Converting queue item to json for build %d for %s", b.GetNumber(), r.GetFullName())
	byteItem, err := json.Marshal(item)
	if err != nil {
		logrus.Errorf("Failed to convert item to json for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)
		return
	}

	logrus.Infof("Publishing item for build %d for %s to queue", b.GetNumber(), r.GetFullName())
	err = queue.Publish("vela", byteItem)
	if err != nil {
		logrus.Errorf("Retrying; Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)
		err = queue.Publish("vela", byteItem)
		if err != nil {
			logrus.Errorf("Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)
			return
		}
	}
}
