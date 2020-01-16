// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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

var baseErr = "unable to process webhook"

// PostWebhook represents the API handler to capture
// a webhook from a source control provider and
// publish it to the configure queue.
func PostWebhook(c *gin.Context) {
	logrus.Info("Webhook received")

	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)

	// process the webhook from the source control provider
	h, r, b, err := source.FromContext(c).ProcessWebhook(c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to parse webhook: %v", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	defer func() {
		// send API call to update the webhook
		err = database.FromContext(c).UpdateHook(h)
		if err != nil {
			logrus.Errorf("unable to update webhook %s/%s: %v", r.GetFullName(), h.GetSourceID(), err)
		}
	}()

	// check if build was parsed from webhook
	if b == nil {
		// typically, this should only happen on a webhook
		// "ping" which gets sent when the webhook is created
		c.JSON(http.StatusOK, "no build to process")

		return
	}

	// check if repo was parsed from webhook
	if r == nil {
		retErr := fmt.Errorf("%s: failed to parse repo from webhook", baseErr)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture parsed repo from webhook
	r, err = database.FromContext(c).GetRepo(r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get repo %s: %v", baseErr, r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// set the RepoID fields
	b.SetRepoID(r.GetID())
	h.SetRepoID(r.GetID())

	// send API call to capture the last hook for the repo
	lastHook, err := database.FromContext(c).GetLastHook(r)
	if err != nil {
		retErr := fmt.Errorf("unable to get last hook for repo %s: %v", r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// set the Number field
	if lastHook != nil {
		h.SetNumber(
			lastHook.GetNumber() + 1,
		)
	}

	// send API call to create the webhook
	err = database.FromContext(c).CreateHook(h)
	if err != nil {
		retErr := fmt.Errorf("unable to create webhook %s/%d: %v", r.GetFullName(), h.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture the created webhook
	h, _ = database.FromContext(c).GetHook(h.GetNumber(), r)

	// check if the repo is active
	if !r.GetActive() {
		retErr := fmt.Errorf("%s: %s is not an active repo", baseErr, r.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// verify the build has a valid event and the repo allows that event type
	if (b.GetEvent() == constants.EventPush && !r.GetAllowPush()) ||
		(b.GetEvent() == constants.EventPull && !r.GetAllowPull()) ||
		(b.GetEvent() == constants.EventTag && !r.GetAllowTag()) ||
		(b.GetEvent() == constants.EventDeploy && !r.GetAllowDeploy()) {
		retErr := fmt.Errorf("%s: %s does not have %s events enabled", baseErr, r.GetFullName(), b.GetEvent())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// check if the repo has a valid owner
	if r.GetUserID() == 0 {
		retErr := fmt.Errorf("%s: %s has no valid owner", baseErr, r.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get owner for %s: %v", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture the last build for the repo
	lastBuild, err := database.FromContext(c).GetLastBuild(r)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get last build for %s: %v", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// update fields in build object
	b.SetNumber(1)
	b.SetParent(b.GetNumber())
	b.SetStatus(constants.StatusPending)
	b.SetCreated(time.Now().UTC().Unix())

	if lastBuild != nil {
		b.SetNumber(
			lastBuild.GetNumber() + 1,
		)
		b.SetParent(lastBuild.GetNumber())
	}

	// populate the build link if a web address is provided
	if len(m.Vela.WebAddress) > 0 {
		b.SetLink(
			fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, r.GetFullName(), b.GetNumber()),
		)
	}

	// variable to store changeset files
	var files []string
	// check if the build event is not pull_request
	if !strings.EqualFold(b.GetEvent(), constants.EventPull) {
		// send API call to capture list of files changed for the commit
		files, err = source.FromContext(c).Changeset(u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %v", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

	// files is empty if the build event is pull_request
	if len(files) == 0 {
		// parse out pull request number from base ref
		//
		// pattern: refs/pull/1/head
		var parts []string
		if strings.HasPrefix(b.GetRef(), "refs/pull/") {
			parts = strings.Split(b.GetRef(), "/")
		}

		// capture number by converting from string
		number, err := strconv.Atoi(parts[2])
		if err != nil {
			// capture number by scanning from string
			_, err := fmt.Sscanf(b.GetRef(), "%s/%s/%d/%s", nil, nil, &number, nil)
			if err != nil {
				retErr := fmt.Errorf("%s: failed to get pull_request number for %s: %v", baseErr, r.GetFullName(), err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}
		}

		// send API call to capture list of files changed for the pull request
		files, err = source.FromContext(c).ChangesetPR(u, r, number)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %v", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).Config(u, r.GetOrg(), r.GetName(), b.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get pipeline configuration for %s: %v", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusNotFound, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// parse and compile the pipeline configuration file
	p, err := compiler.FromContext(c).
		WithBuild(b).
		WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to compile pipeline configuration for %s: %v", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), p, b, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		h.SetStatus(constants.StatusFailure)
		h.SetError(err.Error())

		return
	}

	// send API call to capture the triggered build
	b, _ = database.FromContext(c).GetBuild(b.GetNumber(), r)

	// set the BuildID field
	h.SetBuildID(b.GetID())

	c.JSON(http.StatusOK, b)

	// send API call to set the status on the commit
	err = source.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromContext(c),
		p,
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

	logrus.Infof("Establishing route for build %d for %s", b.GetNumber(), r.GetFullName())

	route, err := queue.Route(&p.Worker)
	if err != nil {
		logrus.Errorf("unable to set route for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		return
	}

	logrus.Infof("Publishing item for build %d for %s to queue", b.GetNumber(), r.GetFullName())

	err = queue.Publish(route, byteItem)
	if err != nil {
		logrus.Errorf("Retrying; Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		err = queue.Publish(route, byteItem)
		if err != nil {
			logrus.Errorf("Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

			return
		}
	}
}
