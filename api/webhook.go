// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var baseErr = "unable to process webhook"

// swagger:operation POST /webhook base PostWebhook
//
// Deliver a webhook to the vela api
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Webhook payload that we expect from the user or VCS
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
// responses:
//   '200':
//     description: Successfully received the webhook
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Malformed webhook payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to receive the webhook
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthenticated
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to receive the webhook
//     schema:
//       "$ref": "#/definitions/Error"

// PostWebhook represents the API handler to capture
// a webhook from a source control provider and
// publish it to the configure queue.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func PostWebhook(c *gin.Context) {
	logrus.Info("webhook received")

	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)

	// duplicate request so we can perform operations on the request body
	//
	// https://golang.org/pkg/net/http/#Request.Clone
	dupRequest := c.Request.Clone(context.TODO())

	// -------------------- Start of TODO: --------------------
	//
	// Remove the below code once http.Request.Clone()
	// actually performs a deep clone.
	//
	// This code is required due to a known bug:
	//
	// * https://github.com/golang/go/issues/36095

	// create buffer for reading request body
	var buf bytes.Buffer

	// read the request body for duplication
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		retErr := fmt.Errorf("unable to read webhook body: %w", err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// add the request body to the original request
	c.Request.Body = io.NopCloser(&buf)

	// add the request body to the duplicate request
	dupRequest.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	//
	// -------------------- End of TODO: --------------------

	// process the webhook from the source control provider
	// comment, number, h, r, b
	webhook, err := scm.FromContext(c).ProcessWebhook(c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to parse webhook: %w", err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// check if the hook should be skipped
	if skip, skipReason := webhook.ShouldSkip(); skip {
		c.JSON(http.StatusOK, fmt.Sprintf("skipping build: %s", skipReason))

		return
	}

	h, r, b := webhook.Hook, webhook.Repo, webhook.Build

	logrus.Debugf("hook generated from SCM: %v", h)
	logrus.Debugf("repo generated from SCM: %v", r)

	if b != nil {
		logrus.Debugf(`build author: %s,
		build branch: %s,
		build commit: %s,
		build ref: %s`,
			b.GetAuthor(), b.GetBranch(), b.GetCommit(), b.GetRef())
	}

	// check if build was parsed from webhook.
	// build will be nil on repository events, but
	// for renaming, we want to continue.
	if b == nil && h.GetEvent() != constants.EventRepository {
		// typically, this should only happen on a webhook
		// "ping" which gets sent when the webhook is created
		c.JSON(http.StatusOK, "no build to process")

		return
	}

	// check if repo was parsed from webhook
	if r == nil {
		retErr := fmt.Errorf("%s: failed to parse repo from webhook", baseErr)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	defer func() {
		// send API call to update the webhook
		err = database.FromContext(c).UpdateHook(h)
		if err != nil {
			logrus.Errorf("unable to update webhook %s/%d: %v", r.GetFullName(), h.GetNumber(), err)
		}
	}()

	if h.GetEvent() == constants.EventRepository {
		switch h.GetEventAction() {
		// if action is rename, go through rename routine
		case constants.ActionRenamed:
			err = renameRepository(h, r, c, m)
			if err != nil {
				util.HandleError(c, http.StatusBadRequest, err)
				h.SetStatus(constants.StatusFailure)
				h.SetError(err.Error())
			}

			c.JSON(http.StatusOK, fmt.Sprintf("no build to process, repository renamed from %s to %s", r.GetPreviousName(), r.GetFullName()))

			return
		// if action is archived, unarchived, or edited, perform edits to relevant repo fields
		case "archived", "unarchived", constants.ActionEdited:
			// send call to get repository from database
			dbRepo, err := database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())
			if err != nil {
				retErr := fmt.Errorf("%s: failed to get repo %s: %w", baseErr, r.GetFullName(), err)
				util.HandleError(c, http.StatusBadRequest, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}

			var retMsg string
			// the only edits to a repo that impact Vela are to these two fields
			if !strings.EqualFold(dbRepo.GetBranch(), r.GetBranch()) {
				retMsg = fmt.Sprintf("no build to process, repository default branch changed from %s to %s", dbRepo.GetBranch(), r.GetBranch())
				dbRepo.SetBranch(r.GetBranch())
			}

			if dbRepo.GetActive() != r.GetActive() {
				retMsg = fmt.Sprintf("no build to process, repository changed active status from %t to %t", dbRepo.GetActive(), r.GetActive())
				dbRepo.SetActive(r.GetActive())
			}

			// update repo object in the database after applying edits
			err = database.FromContext(c).UpdateRepo(dbRepo)
			if err != nil {
				retErr := fmt.Errorf("%s: failed to update repo %s: %w", baseErr, r.GetFullName(), err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}

			c.JSON(http.StatusOK, retMsg)

			return
		// all other repo event actions are skippable
		default:
			c.JSON(http.StatusOK, "no build to process")

			return
		}
	}

	// send API call to capture parsed repo from webhook
	r, err = database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get repo %s: %w", baseErr, r.GetFullName(), err)
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
		retErr := fmt.Errorf("unable to get last hook for repo %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

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
		retErr := fmt.Errorf("unable to create webhook %s/%d: %w", r.GetFullName(), h.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture the created webhook
	h, _ = database.FromContext(c).GetHook(h.GetNumber(), r)

	// verify the webhook from the source control provider
	if c.Value("webhookvalidation").(bool) {
		err = scm.FromContext(c).VerifyWebhook(dupRequest, r)
		if err != nil {
			retErr := fmt.Errorf("unable to verify webhook: %w", err)
			util.HandleError(c, http.StatusUnauthorized, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

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
		(b.GetEvent() == constants.EventComment && !r.GetAllowComment()) ||
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
	logrus.Debugf("capturing owner of repository %s", r.GetFullName())

	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get owner for %s: %w", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.FromContext(c).GetRepoBuildCount(r, filters)
	if err != nil {
		retErr := fmt.Errorf("%s: unable to get count of builds for repo %s", baseErr, r.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	logrus.Debugf("currently %d builds running on repo %s", builds, r.GetFullName())

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= r.GetBuildLimit() {
		retErr := fmt.Errorf("%s: repo %s has exceeded the concurrent build limit of %d", baseErr, r.GetFullName(), r.GetBuildLimit())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// update fields in build object
	logrus.Debugf("updating build number to %d", r.GetCounter())
	b.SetNumber(r.GetCounter())

	logrus.Debugf("updating parent number to %d", b.GetNumber())
	b.SetParent(b.GetNumber())

	logrus.Debug("updating status to pending")
	b.SetStatus(constants.StatusPending)

	// if this is a comment on a pull_request event
	if strings.EqualFold(b.GetEvent(), constants.EventComment) && webhook.PRNumber > 0 {
		commit, branch, baseref, headref, err := scm.FromContext(c).GetPullRequest(u, r, webhook.PRNumber)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get pull request info for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}

		b.SetCommit(commit)
		b.SetBranch(strings.Replace(branch, "refs/heads/", "", -1))
		b.SetBaseRef(baseref)
		b.SetHeadRef(headref)
	}

	// variable to store changeset files
	var files []string
	// check if the build event is not issue_comment or pull_request
	if !strings.EqualFold(b.GetEvent(), constants.EventComment) &&
		!strings.EqualFold(b.GetEvent(), constants.EventPull) {
		// send API call to capture list of files changed for the commit
		files, err = scm.FromContext(c).Changeset(u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

	// check if the build event is a pull_request
	if strings.EqualFold(b.GetEvent(), constants.EventPull) && webhook.PRNumber > 0 {
		// send API call to capture list of files changed for the pull request
		files, err = scm.FromContext(c).ChangesetPR(u, r, webhook.PRNumber)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

	var (
		// variable to store the raw pipeline configuration
		config []byte
		// variable to store executable pipeline
		p *pipeline.Build
		// variable to store pipeline configuration
		pipeline *library.Pipeline
		// variable to control number of times to retry processing pipeline
		retryLimit = 3
		// variable to store the pipeline type for the repository
		pipelineType = r.GetPipelineType()
	)

	// implement a loop to process asynchronous operations with a retry limit
	//
	// Some operations taken during the webhook workflow can lead to race conditions
	// failing to successfully process the request. This logic ensures we attempt our
	// best efforts to handle these cases gracefully.
	for i := 0; i < retryLimit; i++ {
		logrus.Debugf("compilation loop - attempt %d", i+1)
		// check if we're on the first iteration of the loop
		if i > 0 {
			// incrementally sleep in between retries
			time.Sleep(time.Duration(i) * time.Second)
		}

		// send API call to attempt to capture the pipeline
		pipeline, err = database.FromContext(c).GetPipelineForRepo(b.GetCommit(), r)
		if err != nil { // assume the pipeline doesn't exist in the database yet
			// send API call to capture the pipeline configuration file
			config, err = scm.FromContext(c).ConfigBackoff(u, r, b.GetCommit())
			if err != nil {
				retErr := fmt.Errorf("%s: unable to get pipeline configuration for %s: %w", baseErr, r.GetFullName(), err)

				util.HandleError(c, http.StatusNotFound, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}
		} else {
			config = pipeline.GetData()
		}

		// send API call to capture repo for the counter
		r, err = database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())
		if err != nil {
			retErr := fmt.Errorf("%s: unable to get repo %s: %w", baseErr, r.GetFullName(), err)

			// check if the retry limit has been exceeded
			if i < retryLimit-1 {
				logrus.WithError(retErr).Warningf("retrying #%d", i+1)

				// continue to the next iteration of the loop
				continue
			}

			util.HandleError(c, http.StatusBadRequest, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}

		// set the parent equal to the current repo counter
		b.SetParent(r.GetCounter())

		// check if the parent is set to 0
		if b.GetParent() == 0 {
			// parent should be "1" if it's the first build ran
			b.SetParent(1)
		}

		// update the build numbers based off repo counter
		inc := r.GetCounter() + 1
		r.SetCounter(inc)
		b.SetNumber(inc)

		// populate the build link if a web address is provided
		if len(m.Vela.WebAddress) > 0 {
			b.SetLink(
				fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, r.GetFullName(), b.GetNumber()),
			)
		}

		// ensure we use the expected pipeline type when compiling
		//
		// The pipeline type for a repo can change at any time which can break compiling
		// existing pipelines in the system for that repo. To account for this, we update
		// the repo pipeline type to match what was defined for the existing pipeline
		// before compiling. After we're done compiling, we reset the pipeline type.
		if len(pipeline.GetType()) > 0 {
			r.SetPipelineType(pipeline.GetType())
		}

		var compiled *library.Pipeline
		// parse and compile the pipeline configuration file
		p, compiled, err = compiler.FromContext(c).
			Duplicate().
			WithBuild(b).
			WithComment(webhook.Comment).
			WithFiles(files).
			WithMetadata(m).
			WithRepo(r).
			WithUser(u).
			Compile(config)
		if err != nil {
			// format the error message with extra information
			err = fmt.Errorf("unable to compile pipeline configuration for %s: %w", r.GetFullName(), err)

			// log the error for traceability
			logrus.Error(err.Error())

			retErr := fmt.Errorf("%s: %w", baseErr, err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}

		// reset the pipeline type for the repo
		//
		// The pipeline type for a repo can change at any time which can break compiling
		// existing pipelines in the system for that repo. To account for this, we update
		// the repo pipeline type to match what was defined for the existing pipeline
		// before compiling. After we're done compiling, we reset the pipeline type.
		r.SetPipelineType(pipelineType)

		// skip the build if only the init or clone steps are found
		skip := skipEmptyBuild(p)
		if skip != "" {
			// set build to successful status
			b.SetStatus(constants.StatusSkipped)

			// send API call to set the status on the commit
			err = scm.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
			if err != nil {
				logrus.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
			}

			c.JSON(http.StatusOK, skip)

			return
		}

		// check if the pipeline did not already exist in the database
		if pipeline == nil {
			pipeline = compiled
			pipeline.SetRepoID(r.GetID())
			pipeline.SetCommit(b.GetCommit())
			pipeline.SetRef(b.GetRef())

			// send API call to create the pipeline
			err = database.FromContext(c).CreatePipeline(pipeline)
			if err != nil {
				retErr := fmt.Errorf("%s: failed to create pipeline for %s: %w", baseErr, r.GetFullName(), err)

				// check if the retry limit has been exceeded
				if i < retryLimit-1 {
					logrus.WithError(retErr).Warningf("retrying #%d", i+1)

					// continue to the next iteration of the loop
					continue
				}

				util.HandleError(c, http.StatusBadRequest, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}

			// send API call to capture the created pipeline
			pipeline, err = database.FromContext(c).GetPipelineForRepo(pipeline.GetCommit(), r)
			if err != nil {
				//nolint:lll // ignore long line length due to error message
				retErr := fmt.Errorf("%s: failed to get new pipeline %s/%s: %w", baseErr, r.GetFullName(), pipeline.GetCommit(), err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}
		}

		b.SetPipelineID(pipeline.GetID())

		// create the objects from the pipeline in the database
		// TODO:
		// - if a build gets created and something else fails midway,
		//   the next loop will attempt to create the same build,
		//   using the same Number and thus create a constraint
		//   conflict; consider deleting the partially created
		//   build object in the database
		err = planBuild(database.FromContext(c), p, b, r)
		if err != nil {
			retErr := fmt.Errorf("%s: %w", baseErr, err)

			// check if the retry limit has been exceeded
			if i < retryLimit-1 {
				logrus.WithError(retErr).Warningf("retrying #%d", i+1)

				// reset fields set by cleanBuild for retry
				b.SetError("")
				b.SetStatus(constants.StatusPending)
				b.SetFinished(0)

				// continue to the next iteration of the loop
				continue
			}

			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}

		// break the loop because everything was successful
		break
	} // end of retry loop

	// send API call to update repo for ensuring counter is incremented
	err = database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to update repo %s: %w", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// return error if pipeline didn't get populated
	if p == nil {
		retErr := fmt.Errorf("%s: failed to set pipeline for %s: %w", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// return error if build didn't get populated
	if b == nil {
		retErr := fmt.Errorf("%s: failed to set build for %s: %w", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture the triggered build
	b, err = database.FromContext(c).GetBuild(b.GetNumber(), r)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get new build %s/%d: %w", baseErr, r.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// set the BuildID field
	h.SetBuildID(b.GetID())

	c.JSON(http.StatusOK, b)

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromGinContext(c),
		database.FromContext(c),
		p,
		b,
		r,
		u,
	)
}

// publishToQueue is a helper function that creates
// a build item and publishes it to the queue.
func publishToQueue(queue queue.Service, db database.Service, p *pipeline.Build, b *library.Build, r *library.Repo, u *library.User) {
	item := types.ToItem(p, b, r, u)

	logrus.Infof("Converting queue item to json for build %d for %s", b.GetNumber(), r.GetFullName())

	byteItem, err := json.Marshal(item)
	if err != nil {
		logrus.Errorf("Failed to convert item to json for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		cleanBuild(db, b, nil, nil)

		return
	}

	logrus.Infof("Establishing route for build %d for %s", b.GetNumber(), r.GetFullName())

	route, err := queue.Route(&p.Worker)
	if err != nil {
		logrus.Errorf("unable to set route for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		cleanBuild(db, b, nil, nil)

		return
	}

	logrus.Infof("Publishing item for build %d for %s to queue %s", b.GetNumber(), r.GetFullName(), route)

	err = queue.Push(context.Background(), route, byteItem)
	if err != nil {
		logrus.Errorf("Retrying; Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		err = queue.Push(context.Background(), route, byteItem)
		if err != nil {
			logrus.Errorf("Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

			// error out the build
			cleanBuild(db, b, nil, nil)

			return
		}
	}

	// update fields in build object
	b.SetEnqueued(time.Now().UTC().Unix())

	// update the build in the db to reflect the time it was enqueued
	err = db.UpdateBuild(b)
	if err != nil {
		logrus.Errorf("Failed to update build %d during publish to queue for %s: %v", b.GetNumber(), r.GetFullName(), err)
	}
}

// renameRepository is a helper function that takes the old name of the repo,
// queries the database for the repo that matches that name and org, and updates
// that repo to its new name in order to preserve it. It also updates the secrets
// associated with that repo.
func renameRepository(h *library.Hook, r *library.Repo, c *gin.Context, m *types.Metadata) error {
	logrus.Debugf("renaming repository from %s to %s", r.GetPreviousName(), r.GetName())
	// get the old name of the repo
	previousName := r.GetPreviousName()
	// get the repo from the database that matches the old name
	dbR, err := database.FromContext(c).GetRepoForOrg(r.GetOrg(), previousName)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get repo %s/%s from database", baseErr, r.GetOrg(), previousName)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return retErr
	}

	// update the repo name information
	dbR.SetName(r.GetName())
	dbR.SetFullName(r.GetFullName())
	dbR.SetClone(r.GetClone())
	dbR.SetLink(r.GetLink())
	dbR.SetPreviousName(previousName)

	// update the repo in the database
	err = database.FromContext(c).UpdateRepo(dbR)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to update repo %s/%s in database", baseErr, r.GetOrg(), previousName)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return retErr
	}

	// update hook object which will be added to DB upon reaching deferred function in PostWebhook
	h.SetRepoID(r.GetID())

	// send API call to capture the last hook for the repo
	lastHook, err := database.FromContext(c).GetLastHook(dbR)
	if err != nil {
		retErr := fmt.Errorf("unable to get last hook for repo %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return retErr
	}

	// set the Number field
	if lastHook != nil {
		h.SetNumber(
			lastHook.GetNumber() + 1,
		)
	}

	// get total number of secrets associated with repository
	t, err := database.FromContext(c).GetTypeSecretCount(constants.SecretRepo, r.GetOrg(), previousName, []string{})
	if err != nil {
		return fmt.Errorf("unable to get secret count for repo %s/%s: %w", r.GetOrg(), previousName, err)
	}

	secrets := []*library.Secret{}
	page := 1
	// capture all secrets belonging to certain repo in database
	for repoSecrets := int64(0); repoSecrets < t; repoSecrets += 100 {
		s, err := database.FromContext(c).GetTypeSecretList(constants.SecretRepo, r.GetOrg(), previousName, page, 100, []string{})
		if err != nil {
			return fmt.Errorf("unable to get secret list for repo %s/%s: %w", r.GetOrg(), previousName, err)
		}

		secrets = append(secrets, s...)

		page++
	}

	// update secrets to point to the new repository name
	for _, secret := range secrets {
		secret.SetRepo(r.GetName())

		err = database.FromContext(c).UpdateSecret(secret)
		if err != nil {
			return fmt.Errorf("unable to update secret for repo %s/%s: %w", r.GetOrg(), previousName, err)
		}
	}

	// get total number of builds associated with repository
	t, err = database.FromContext(c).GetRepoBuildCount(dbR, nil)
	if err != nil {
		return fmt.Errorf("unable to get build count for repo %s: %w", dbR.GetFullName(), err)
	}

	builds := []*library.Build{}
	page = 1
	// capture all builds belonging to repo in database
	for build := int64(0); build < t; build += 100 {
		b, _, err := database.FromContext(c).GetRepoBuildList(dbR, nil, time.Now().Unix(), 0, page, 100)
		if err != nil {
			return fmt.Errorf("unable to get build list for repo %s: %w", dbR.GetFullName(), err)
		}

		builds = append(builds, b...)

		page++
	}

	// update build link to route to proper repo name
	for _, build := range builds {
		build.SetLink(
			fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, dbR.GetFullName(), build.GetNumber()),
		)

		err = database.FromContext(c).UpdateBuild(build)
		if err != nil {
			return fmt.Errorf("unable to update build for repo %s: %w", dbR.GetFullName(), err)
		}
	}

	return nil
}
