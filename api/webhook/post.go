// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/build"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	ctx := c.Request.Context()

	// duplicate request so we can perform operations on the request body
	//
	// https://golang.org/pkg/net/http/#Request.Clone
	dupRequest := c.Request.Clone(ctx)

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
	//
	// populate build, hook, repo resources as well as PR Number / PR Comment if necessary
	webhook, err := scm.FromContext(c).ProcessWebhook(ctx, c.Request)
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

	// if event is repository event, handle separately and return
	if strings.EqualFold(h.GetEvent(), constants.EventRepository) {
		r, err = handleRepositoryEvent(ctx, c, m, h, r)
		if err != nil {
			util.HandleError(c, http.StatusInternalServerError, err)
			return
		}

		// if there were actual changes to the repo (database call populated ID field), return the repo object
		if r.GetID() != 0 {
			c.JSON(http.StatusOK, r)
			return
		}

		c.JSON(http.StatusOK, "handled repository event, no build to process")

		return
	}

	// check if build was parsed from webhook.
	if b == nil {
		// typically, this should only happen on a webhook
		// "ping" which gets sent when the webhook is created
		c.JSON(http.StatusOK, "no build to process")

		return
	}

	logrus.Debugf(`build author: %s,
		build branch: %s,
		build commit: %s,
		build ref: %s`,
		b.GetAuthor(), b.GetBranch(), b.GetCommit(), b.GetRef())

	// check if repo was parsed from webhook
	if r == nil {
		retErr := fmt.Errorf("%s: failed to parse repo from webhook", baseErr)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	defer func() {
		// send API call to update the webhook
		_, err = database.FromContext(c).UpdateHook(ctx, h)
		if err != nil {
			logrus.Errorf("unable to update webhook %s/%d: %v", r.GetFullName(), h.GetNumber(), err)
		}
	}()

	// send API call to capture parsed repo from webhook
	repo, err := database.FromContext(c).GetRepoForOrg(ctx, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get repo %s: %w", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// set the RepoID fields
	b.SetRepoID(repo.GetID())
	h.SetRepoID(repo.GetID())

	// send API call to capture the last hook for the repo
	lastHook, err := database.FromContext(c).LastHookForRepo(ctx, repo)
	if err != nil {
		retErr := fmt.Errorf("unable to get last hook for repo %s: %w", repo.GetFullName(), err)
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
	h, err = database.FromContext(c).CreateHook(ctx, h)
	if err != nil {
		retErr := fmt.Errorf("unable to create webhook %s/%d: %w", repo.GetFullName(), h.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// verify the webhook from the source control provider
	if c.Value("webhookvalidation").(bool) {
		err = scm.FromContext(c).VerifyWebhook(ctx, dupRequest, repo)
		if err != nil {
			retErr := fmt.Errorf("unable to verify webhook: %w", err)
			util.HandleError(c, http.StatusUnauthorized, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

	// check if the repo is active
	if !repo.GetActive() {
		retErr := fmt.Errorf("%s: %s is not an active repo", baseErr, repo.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// verify the build has a valid event and the repo allows that event type
	if !repo.EventAllowed(b.GetEvent(), b.GetEventAction()) {
		var actionErr string
		if len(b.GetEventAction()) > 0 {
			actionErr = ":" + b.GetEventAction()
		}

		retErr := fmt.Errorf("%s: %s does not have %s%s events enabled", baseErr, repo.GetFullName(), b.GetEvent(), actionErr)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusSkipped)
		h.SetError(retErr.Error())

		return
	}

	// check if the repo has a valid owner
	if repo.GetUserID() == 0 {
		retErr := fmt.Errorf("%s: %s has no valid owner", baseErr, repo.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture repo owner
	logrus.Debugf("capturing owner of repository %s", repo.GetFullName())

	u, err := database.FromContext(c).GetUser(ctx, repo.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get owner for %s: %w", baseErr, repo.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// confirm current repo owner has at least write access to repo (needed for status update later)
	_, err = scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to publish build to queue: repository owner %s no longer has write access to repository %s", u.GetName(), r.GetFullName())
		util.HandleError(c, http.StatusUnauthorized, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.FromContext(c).CountBuildsForRepo(ctx, repo, filters)
	if err != nil {
		retErr := fmt.Errorf("%s: unable to get count of builds for repo %s", baseErr, repo.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	logrus.Debugf("currently %d builds running on repo %s", builds, repo.GetFullName())

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= repo.GetBuildLimit() {
		retErr := fmt.Errorf("%s: repo %s has exceeded the concurrent build limit of %d", baseErr, repo.GetFullName(), repo.GetBuildLimit())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// update fields in build object
	logrus.Debugf("updating build number to %d", repo.GetCounter())
	b.SetNumber(repo.GetCounter())

	logrus.Debug("updating status to pending")
	b.SetStatus(constants.StatusPending)

	// if the event is issue_comment and the issue is a pull request,
	// call SCM for more data not provided in webhook payload
	if strings.EqualFold(b.GetEvent(), constants.EventComment) && webhook.PullRequest.Number > 0 {
		commit, branch, baseref, headref, err := scm.FromContext(c).GetPullRequest(ctx, u, repo, webhook.PullRequest.Number)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get pull request info for %s: %w", baseErr, repo.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}

		b.SetCommit(commit)
		b.SetBranch(strings.ReplaceAll(branch, "refs/heads/", ""))
		b.SetBaseRef(baseref)
		b.SetHeadRef(headref)
	}

	// variable to store changeset files
	var files []string

	// check if the build event is not issue_comment or pull_request
	if !strings.EqualFold(b.GetEvent(), constants.EventComment) &&
		!strings.EqualFold(b.GetEvent(), constants.EventPull) {
		// send API call to capture list of files changed for the commit
		files, err = scm.FromContext(c).Changeset(ctx, u, repo, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %w", baseErr, repo.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

	// check if the build event is a pull_request
	if strings.EqualFold(b.GetEvent(), constants.EventPull) && webhook.PullRequest.Number > 0 {
		// send API call to capture list of files changed for the pull request
		files, err = scm.FromContext(c).ChangesetPR(ctx, u, repo, webhook.PullRequest.Number)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %w", baseErr, repo.GetFullName(), err)
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
		pipelineType = repo.GetPipelineType()
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

		// send database call to attempt to capture the pipeline if we already processed it before
		pipeline, err = database.FromContext(c).GetPipelineForRepo(ctx, b.GetCommit(), repo)
		if err != nil { // assume the pipeline doesn't exist in the database yet
			// send API call to capture the pipeline configuration file
			config, err = scm.FromContext(c).ConfigBackoff(ctx, u, repo, b.GetCommit())
			if err != nil {
				retErr := fmt.Errorf("%s: unable to get pipeline configuration for %s: %w", baseErr, repo.GetFullName(), err)

				util.HandleError(c, http.StatusNotFound, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}
		} else {
			config = pipeline.GetData()
		}

		// send API call to capture repo for the counter (grabbing repo again to ensure counter is correct)
		repo, err = database.FromContext(c).GetRepoForOrg(ctx, repo.GetOrg(), repo.GetName())
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

		// update DB record of repo (repo) with any changes captured from webhook payload (r)
		repo.SetTopics(r.GetTopics())
		repo.SetBranch(r.GetBranch())

		// update the build numbers based off repo counter
		inc := repo.GetCounter() + 1
		repo.SetCounter(inc)
		b.SetNumber(inc)

		// populate the build link if a web address is provided
		if len(m.Vela.WebAddress) > 0 {
			b.SetLink(
				fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, repo.GetFullName(), b.GetNumber()),
			)
		}

		// ensure we use the expected pipeline type when compiling
		//
		// The pipeline type for a repo can change at any time which can break compiling
		// existing pipelines in the system for that repo. To account for this, we update
		// the repo pipeline type to match what was defined for the existing pipeline
		// before compiling. After we're done compiling, we reset the pipeline type.
		if len(pipeline.GetType()) > 0 {
			repo.SetPipelineType(pipeline.GetType())
		}

		var compiled *library.Pipeline
		// parse and compile the pipeline configuration file
		p, compiled, err = compiler.FromContext(c).
			Duplicate().
			WithBuild(b).
			WithComment(webhook.PullRequest.Comment).
			WithCommit(b.GetCommit()).
			WithFiles(files).
			WithMetadata(m).
			WithRepo(repo).
			WithUser(u).
			Compile(config)
		if err != nil {
			// format the error message with extra information
			err = fmt.Errorf("unable to compile pipeline configuration for %s: %w", repo.GetFullName(), err)

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
		repo.SetPipelineType(pipelineType)

		// skip the build if pipeline compiled to only the init and clone steps
		skip := build.SkipEmptyBuild(p)
		if skip != "" {
			// set build to successful status
			b.SetStatus(constants.StatusSkipped)

			// set hook status and message
			h.SetStatus(constants.StatusSkipped)
			h.SetError(skip)

			// send API call to set the status on the commit
			err = scm.FromContext(c).Status(ctx, u, b, repo.GetOrg(), repo.GetName())
			if err != nil {
				logrus.Errorf("unable to set commit status for %s/%d: %v", repo.GetFullName(), b.GetNumber(), err)
			}

			c.JSON(http.StatusOK, skip)

			return
		}

		// check if the pipeline did not already exist in the database
		if pipeline == nil {
			pipeline = compiled
			pipeline.SetRepoID(repo.GetID())
			pipeline.SetCommit(b.GetCommit())
			pipeline.SetRef(b.GetRef())

			// send API call to create the pipeline
			pipeline, err = database.FromContext(c).CreatePipeline(ctx, pipeline)
			if err != nil {
				retErr := fmt.Errorf("%s: failed to create pipeline for %s: %w", baseErr, repo.GetFullName(), err)

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
		}

		b.SetPipelineID(pipeline.GetID())

		// create the objects from the pipeline in the database
		// TODO:
		// - if a build gets created and something else fails midway,
		//   the next loop will attempt to create the same build,
		//   using the same Number and thus create a constraint
		//   conflict; consider deleting the partially created
		//   build object in the database
		err = build.PlanBuild(ctx, database.FromContext(c), p, b, repo)
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
	repo, err = database.FromContext(c).UpdateRepo(ctx, repo)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to update repo %s: %w", baseErr, repo.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// return error if pipeline didn't get populated
	if p == nil {
		retErr := fmt.Errorf("%s: failed to set pipeline for %s: %w", baseErr, repo.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// return error if build didn't get populated
	if b == nil {
		retErr := fmt.Errorf("%s: failed to set build for %s: %w", baseErr, repo.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// send API call to capture the triggered build
	b, err = database.FromContext(c).GetBuildForRepo(ctx, repo, b.GetNumber())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get new build %s/%d: %w", baseErr, repo.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	// set the BuildID field
	h.SetBuildID(b.GetID())

	// if event is deployment, update the deployment record to include this build
	if b.GetEvent() == constants.EventDeploy {
		builds := []*library.Build{}
		builds = append(builds, b)

		d, err := database.FromContext(c).GetDeploymentForRepo(c, repo, webhook.Deployment.GetNumber())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				deployment := webhook.Deployment

				deployment.SetRepoID(repo.GetID())
				deployment.SetBuilds(builds)

				_, err := database.FromContext(c).CreateDeployment(c, deployment)
				if err != nil {
					retErr := fmt.Errorf("%s: failed to create deployment %s/%d: %w", baseErr, repo.GetFullName(), deployment.GetNumber(), err)
					util.HandleError(c, http.StatusInternalServerError, retErr)

					h.SetStatus(constants.StatusFailure)
					h.SetError(retErr.Error())

					return
				}
			} else {
				retErr := fmt.Errorf("%s: failed to get deployment %s/%d: %w", baseErr, repo.GetFullName(), webhook.Deployment.GetNumber(), err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}
		} else {
			d.SetBuilds(builds)
			_, err := database.FromContext(c).UpdateDeployment(d)
			if err != nil {
				retErr := fmt.Errorf("%s: failed to update deployment %s/%d: %w", baseErr, repo.GetFullName(), d.GetNumber(), err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}
		}
	}

	c.JSON(http.StatusOK, b)

	// determine queue route
	route, err := queue.FromContext(c).Route(&p.Worker)
	if err != nil {
		logrus.Errorf("unable to set route for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		build.CleanBuild(ctx, database.FromContext(c), b, nil, nil, err)

		return
	}

	// temporarily set host to the route before it gets picked up by a worker
	b.SetHost(route)

	// publish the pipeline.Build to the build_executables table to be requested by a worker
	err = build.PublishBuildExecutable(ctx, database.FromContext(c), p, b)
	if err != nil {
		retErr := fmt.Errorf("unable to publish build executable for %s/%d: %w", repo.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// if the webhook was from a Pull event from a forked repository, verify it is allowed to run
	if webhook.PullRequest.IsFromFork {
		switch repo.GetApproveBuild() {
		case constants.ApproveForkAlways:
			err = gatekeepBuild(c, b, repo, u)
			if err != nil {
				util.HandleError(c, http.StatusInternalServerError, err)
			}

			return
		case constants.ApproveForkNoWrite:
			// determine if build sender has write access to parent repo. If not, this call will result in an error
			_, err = scm.FromContext(c).RepoAccess(ctx, b.GetSender(), u.GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				err = gatekeepBuild(c, b, repo, u)
				if err != nil {
					util.HandleError(c, http.StatusInternalServerError, err)
				}

				return
			}

			fallthrough
		case constants.ApproveNever:
			fallthrough
		default:
			logrus.Debugf("fork PR build %s/%d automatically running without approval", repo.GetFullName(), b.GetNumber())
		}
	}

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(ctx, u, b, repo.GetOrg(), repo.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for %s/%d: %v", repo.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go build.PublishToQueue(
		ctx,
		queue.FromGinContext(c),
		database.FromContext(c),
		b,
		repo,
		u,
		route,
	)

	if build.ShouldAutoCancel(p.Metadata.AutoCancel, b, repo.GetBranch()) {
		// fetch pending and running builds
		rBs, err := database.FromContext(c).ListPendingAndRunningBuildsForRepo(c, repo)
		if err != nil {
			logrus.Errorf("unable to fetch pending and running builds for %s: %v", repo.GetFullName(), err)
		}

		for _, rB := range rBs {
			// call auto cancel routine
			canceled, err := build.AutoCancel(c, b, rB, repo, p.Metadata.AutoCancel)
			if err != nil {
				// continue cancel loop if error, but log based on type of error
				if canceled {
					logrus.Errorf("unable to update canceled build error message: %v", err)
				} else {
					logrus.Errorf("unable to cancel running build: %v", err)
				}
			}
		}
	}
}

// handleRepositoryEvent is a helper function that processes repository events from the SCM and updates
// the database resources with any relevant changes resulting from the event, such as name changes, transfers, etc.
func handleRepositoryEvent(ctx context.Context, c *gin.Context, m *types.Metadata, h *library.Hook, r *library.Repo) (*library.Repo, error) {
	logrus.Debugf("webhook is repository event, making necessary updates to repo %s", r.GetFullName())

	defer func() {
		// send API call to update the webhook
		_, err := database.FromContext(c).CreateHook(ctx, h)
		if err != nil {
			logrus.Errorf("unable to create webhook %s/%d: %v", r.GetFullName(), h.GetNumber(), err)
		}
	}()

	switch h.GetEventAction() {
	// if action is renamed or transferred, go through rename routine
	case constants.ActionRenamed, constants.ActionTransferred:
		r, err := RenameRepository(ctx, h, r, c, m)
		if err != nil {
			h.SetStatus(constants.StatusFailure)
			h.SetError(err.Error())

			return nil, err
		}

		return r, nil
	// if action is archived, unarchived, or edited, perform edits to relevant repo fields
	case "archived", "unarchived", constants.ActionEdited:
		logrus.Debugf("repository action %s for %s", h.GetEventAction(), r.GetFullName())
		// send call to get repository from database
		dbRepo, err := database.FromContext(c).GetRepoForOrg(ctx, r.GetOrg(), r.GetName())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get repo %s: %w", baseErr, r.GetFullName(), err)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return nil, retErr
		}

		// send API call to capture the last hook for the repo
		lastHook, err := database.FromContext(c).LastHookForRepo(ctx, dbRepo)
		if err != nil {
			retErr := fmt.Errorf("unable to get last hook for repo %s: %w", r.GetFullName(), err)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return nil, retErr
		}

		// set the Number field
		if lastHook != nil {
			h.SetNumber(
				lastHook.GetNumber() + 1,
			)
		}

		h.SetRepoID(dbRepo.GetID())

		// the only edits to a repo that impact Vela are to these three fields
		if !strings.EqualFold(dbRepo.GetBranch(), r.GetBranch()) {
			dbRepo.SetBranch(r.GetBranch())
		}

		if dbRepo.GetActive() != r.GetActive() {
			dbRepo.SetActive(r.GetActive())
		}

		if !reflect.DeepEqual(dbRepo.GetTopics(), r.GetTopics()) {
			dbRepo.SetTopics(r.GetTopics())
		}

		// update repo object in the database after applying edits
		dbRepo, err = database.FromContext(c).UpdateRepo(ctx, dbRepo)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to update repo %s: %w", baseErr, r.GetFullName(), err)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return nil, err
		}

		return dbRepo, nil
	// all other repo event actions are skippable
	default:
		return r, nil
	}
}

// RenameRepository is a helper function that takes the old name of the repo,
// queries the database for the repo that matches that name and org, and updates
// that repo to its new name in order to preserve it. It also updates the secrets
// associated with that repo as well as build links for the UI.
func RenameRepository(ctx context.Context, h *library.Hook, r *library.Repo, c *gin.Context, m *types.Metadata) (*library.Repo, error) {
	logrus.Infof("renaming repository from %s to %s", r.GetPreviousName(), r.GetName())

	// get any matching hook with the repo's unique webhook ID in the SCM
	hook, err := database.FromContext(c).GetHookByWebhookID(ctx, h.GetWebhookID())
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get hook with webhook ID %d from database", baseErr, h.GetWebhookID())
	}

	// get the repo from the database using repo id of matching hook
	dbR, err := database.FromContext(c).GetRepo(ctx, hook.GetRepoID())
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get repo %d from database", baseErr, hook.GetRepoID())
	}

	// update hook object which will be added to DB upon reaching deferred function in PostWebhook
	h.SetRepoID(r.GetID())

	// send API call to capture the last hook for the repo
	lastHook, err := database.FromContext(c).LastHookForRepo(ctx, dbR)
	if err != nil {
		retErr := fmt.Errorf("unable to get last hook for repo %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return nil, retErr
	}

	// set the Number field
	if lastHook != nil {
		h.SetNumber(
			lastHook.GetNumber() + 1,
		)
	}

	// get total number of secrets associated with repository
	t, err := database.FromContext(c).CountSecretsForRepo(ctx, dbR, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unable to get secret count for repo %s/%s: %w", dbR.GetOrg(), dbR.GetName(), err)
	}

	secrets := []*library.Secret{}
	page := 1
	// capture all secrets belonging to certain repo in database
	for repoSecrets := int64(0); repoSecrets < t; repoSecrets += 100 {
		s, _, err := database.FromContext(c).ListSecretsForRepo(ctx, dbR, map[string]interface{}{}, page, 100)
		if err != nil {
			return nil, fmt.Errorf("unable to get secret list for repo %s/%s: %w", dbR.GetOrg(), dbR.GetName(), err)
		}

		secrets = append(secrets, s...)

		page++
	}

	// update secrets to point to the new repository name
	for _, secret := range secrets {
		secret.SetOrg(r.GetOrg())
		secret.SetRepo(r.GetName())

		_, err = database.FromContext(c).UpdateSecret(ctx, secret)
		if err != nil {
			return nil, fmt.Errorf("unable to update secret for repo %s/%s: %w", dbR.GetOrg(), dbR.GetName(), err)
		}
	}

	// get total number of builds associated with repository
	t, err = database.FromContext(c).CountBuildsForRepo(ctx, dbR, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get build count for repo %s: %w", dbR.GetFullName(), err)
	}

	builds := []*library.Build{}
	page = 1
	// capture all builds belonging to repo in database
	for build := int64(0); build < t; build += 100 {
		b, _, err := database.FromContext(c).ListBuildsForRepo(ctx, dbR, nil, time.Now().Unix(), 0, page, 100)
		if err != nil {
			return nil, fmt.Errorf("unable to get build list for repo %s: %w", dbR.GetFullName(), err)
		}

		builds = append(builds, b...)

		page++
	}

	// update build link to route to proper repo name
	for _, build := range builds {
		build.SetLink(
			fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, r.GetFullName(), build.GetNumber()),
		)

		_, err = database.FromContext(c).UpdateBuild(ctx, build)
		if err != nil {
			return nil, fmt.Errorf("unable to update build for repo %s: %w", dbR.GetFullName(), err)
		}
	}

	// update the repo name information
	dbR.SetName(r.GetName())
	dbR.SetOrg(r.GetOrg())
	dbR.SetFullName(r.GetFullName())
	dbR.SetClone(r.GetClone())
	dbR.SetLink(r.GetLink())
	dbR.SetPreviousName(r.GetPreviousName())

	// update the repo in the database
	dbR, err = database.FromContext(c).UpdateRepo(ctx, dbR)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to update repo %s/%s in database", baseErr, dbR.GetOrg(), dbR.GetName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return nil, retErr
	}

	return dbR, nil
}

// gatekeepBuild is a helper function that will set the status of a build to 'pending approval' and
// send a status update to the SCM.
func gatekeepBuild(c *gin.Context, b *library.Build, r *library.Repo, u *library.User) error {
	logrus.Debugf("fork PR build %s/%d waiting for approval", r.GetFullName(), b.GetNumber())
	b.SetStatus(constants.StatusPendingApproval)

	_, err := database.FromContext(c).UpdateBuild(c, b)
	if err != nil {
		return fmt.Errorf("unable to update build for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(c, u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	return nil
}
