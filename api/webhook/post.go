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
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
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
	m := c.MustGet("metadata").(*internal.Metadata)
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
	if !repo.GetAllowEvents().Allowed(b.GetEvent(), b.GetEventAction()) {
		var actionErr string
		if len(b.GetEventAction()) > 0 {
			actionErr = ":" + b.GetEventAction()
		}

		retErr := fmt.Errorf("%s: %s does not have %s%s event enabled", baseErr, repo.GetFullName(), b.GetEvent(), actionErr)
		util.HandleError(c, http.StatusBadRequest, retErr)

		h.SetStatus(constants.StatusSkipped)
		h.SetError(retErr.Error())

		return
	}

	var (
		prComment string
		prLabels  []string
	)

	if strings.EqualFold(b.GetEvent(), constants.EventComment) {
		prComment = webhook.PullRequest.Comment
	}

	if strings.EqualFold(b.GetEvent(), constants.EventPull) {
		prLabels = webhook.PullRequest.Labels
	}

	// construct CompileAndPublishConfig
	config := build.CompileAndPublishConfig{
		Build:    b,
		Repo:     repo,
		Metadata: m,
		BaseErr:  baseErr,
		Source:   "webhook",
		Comment:  prComment,
		Labels:   prLabels,
		Retries:  3,
	}

	// generate the queue item
	p, item, err := build.CompileAndPublish(
		c,
		config,
		database.FromContext(c),
		scm.FromContext(c),
		compiler.FromContext(c),
		queue.FromContext(c),
	)

	// capture the build, repo, and user from the items
	b, repo = item.Build, item.Repo

	// set hook build_id to the generated build id
	h.SetBuildID(b.GetID())

	// check if build was skipped
	if err != nil && strings.EqualFold(b.GetStatus(), constants.StatusSkipped) {
		h.SetStatus(constants.StatusSkipped)
		h.SetError(err.Error())

		c.JSON(http.StatusOK, err.Error())

		return
	}

	if err != nil {
		h.SetStatus(constants.StatusFailure)
		h.SetError(err.Error())

		return
	}

	// if event is deployment, update the deployment record to include this build
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		d, err := database.FromContext(c).GetDeploymentForRepo(c, repo, webhook.Deployment.GetNumber())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				deployment := webhook.Deployment

				deployment.SetRepoID(repo.GetID())
				deployment.SetBuilds([]*library.Build{b})

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
			build := append(d.GetBuilds(), b)
			d.SetBuilds(build)
			_, err := database.FromContext(c).UpdateDeployment(ctx, d)
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

	// regardless of whether the build is published to queue, we want to attempt to auto-cancel if no errors
	defer func() {
		if err == nil && build.ShouldAutoCancel(p.Metadata.AutoCancel, b, repo.GetBranch()) {
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
	}()

	// if the webhook was from a Pull event from a forked repository, verify it is allowed to run
	if webhook.PullRequest.IsFromFork {
		logrus.Tracef("inside %s workflow for fork PR build %s/%d", repo.GetApproveBuild(), r.GetFullName(), b.GetNumber())

		switch repo.GetApproveBuild() {
		case constants.ApproveForkAlways:
			err = gatekeepBuild(c, b, repo)
			if err != nil {
				util.HandleError(c, http.StatusInternalServerError, err)
			}

			return
		case constants.ApproveForkNoWrite:
			// determine if build sender has write access to parent repo. If not, this call will result in an error
			_, err = scm.FromContext(c).RepoAccess(ctx, b.GetSender(), r.GetOwner().GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				err = gatekeepBuild(c, b, repo)
				if err != nil {
					util.HandleError(c, http.StatusInternalServerError, err)
				}

				return
			}

			logrus.Debugf("fork PR build %s/%d automatically running without approval", repo.GetFullName(), b.GetNumber())
		case constants.ApproveOnce:
			// determine if build sender is in the contributors list for the repo
			//
			// NOTE: this call is cumbersome for repos with lots of contributors. Potential TODO: improve this if
			// GitHub adds a single-contributor API endpoint.
			contributor, err := scm.FromContext(c).RepoContributor(ctx, r.GetOwner(), b.GetSender(), r.GetOrg(), r.GetName())
			if err != nil {
				util.HandleError(c, http.StatusInternalServerError, err)
			}

			if !contributor {
				err = gatekeepBuild(c, b, repo)
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
	err = scm.FromContext(c).Status(ctx, repo.GetOwner(), b, repo.GetOrg(), repo.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for %s/%d: %v", repo.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go build.Enqueue(
		ctx,
		queue.FromGinContext(c),
		database.FromContext(c),
		item,
		b.GetHost(),
	)
}

// handleRepositoryEvent is a helper function that processes repository events from the SCM and updates
// the database resources with any relevant changes resulting from the event, such as name changes, transfers, etc.
func handleRepositoryEvent(ctx context.Context, c *gin.Context, m *internal.Metadata, h *library.Hook, r *types.Repo) (*types.Repo, error) {
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
func RenameRepository(ctx context.Context, h *library.Hook, r *types.Repo, c *gin.Context, m *internal.Metadata) (*types.Repo, error) {
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
func gatekeepBuild(c *gin.Context, b *library.Build, r *types.Repo) error {
	logrus.Debugf("fork PR build %s/%d waiting for approval", r.GetFullName(), b.GetNumber())
	b.SetStatus(constants.StatusPendingApproval)

	_, err := database.FromContext(c).UpdateBuild(c, b)
	if err != nil {
		return fmt.Errorf("unable to update build for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// update the build components to pending approval status
	err = build.UpdateComponentStatuses(c, b, constants.StatusPendingApproval)
	if err != nil {
		return fmt.Errorf("unable to update build components for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(c, r.GetOwner(), b, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	return nil
}
