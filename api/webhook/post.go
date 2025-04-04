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
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/api/build"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

var baseErr = "unable to process webhook"

// swagger:operation POST /webhook base PostWebhook
//
// Deliver a webhook to the Vela API
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
//     description: Successfully received the webhook but build was skipped
//     schema:
//       type: string
//   '201':
//     description: Successfully created the build from webhook
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '429':
//     description: Concurrent build limit reached for repository
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// PostWebhook represents the API handler to capture
// a webhook from a source control provider and
// publish it to the configure queue.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func PostWebhook(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	db := database.FromContext(c)
	ctx := c.Request.Context()

	l.Debug("webhook received")

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

	if webhook.Installation != nil {
		l.Trace("verifying GitHub App webhook")

		if c.Value("webhookvalidation").(bool) {
			webhookSecret := c.MustGet("app-webhook-secret").(string)

			err = scm.FromContext(c).VerifyWebhook(ctx, dupRequest, []byte(webhookSecret))
			if err != nil {
				retErr := fmt.Errorf("unable to verify webhook: %w", err)
				util.HandleError(c, http.StatusUnauthorized, retErr)

				return
			}
		}

		err = scm.FromContext(c).ProcessInstallation(ctx, c.Request, webhook, db)
		if err != nil {
			retErr := fmt.Errorf("unable to process installation: %w", err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		c.JSON(http.StatusOK, "installation processed successfully")

		return
	}

	// check if the hook should be skipped
	if skip, skipReason := webhook.ShouldSkip(); skip {
		c.JSON(http.StatusOK, fmt.Sprintf("skipping build: %s", skipReason))

		return
	}

	h, r, b := webhook.Hook, webhook.Repo, webhook.Build

	// check if repo was parsed from webhook
	if r == nil {
		retErr := fmt.Errorf("%s: failed to parse repo from webhook", baseErr)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture parsed repo from webhook
	repo, err := database.FromContext(c).GetRepoForOrg(ctx, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get repo %s: %w", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

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

	l.Debugf("hook generated from SCM: %v", h)
	l.Debugf("repo generated from SCM: %v", r)

	// check if build was parsed from webhook.
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

	var repo *types.Repo

	if h.GetEvent() == constants.EventRepository && (h.GetEventAction() == constants.ActionRenamed || h.GetEventAction() == constants.ActionTransferred) {
		// get any matching hook with the repo's unique webhook ID in the SCM
		hook, err := db.GetHookByWebhookID(ctx, h.GetWebhookID())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get hook by webhook id for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// get the repo from the database using repo id of matching hook
		repo, err = db.GetRepo(ctx, hook.GetRepo().GetID())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get repo by id: %w", baseErr, err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	} else {
		repo, err = database.FromContext(c).GetRepoForOrg(ctx, r.GetOrg(), r.GetName())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get repo %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	// verify the webhook from the source control provider using DB repo hash
	if c.Value("webhookvalidation").(bool) {
		l.WithFields(logrus.Fields{
			"org":  r.GetOrg(),
			"repo": r.GetName(),
		}).Tracef("verifying GitHub webhook for %s", r.GetFullName())

		err = scm.FromContext(c).VerifyWebhook(ctx, dupRequest, []byte(repo.GetHash()))
		if err != nil {
			retErr := fmt.Errorf("unable to verify webhook: %w", err)
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}

	// if event is repository event, handle separately and return
	if strings.EqualFold(h.GetEvent(), constants.EventRepository) {
		r, err = handleRepositoryEvent(ctx, l, db, m, h, r, repo)
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

	l.Debugf(`build author: %s,
		build branch: %s,
		build commit: %s,
		build ref: %s`,
		b.GetAuthor(), b.GetBranch(), b.GetCommit(), b.GetRef())

	defer func() {
		// send API call to update the webhook
		//
		//nolint:contextcheck // false positive
		_, err = database.FromContext(c).UpdateHook(ctx, h)
		if err != nil {
			l.Errorf("unable to update webhook %s/%d: %v", r.GetFullName(), h.GetNumber(), err)
		}

		l.WithFields(logrus.Fields{
			"hook":    h.GetNumber(),
			"hook_id": h.GetID(),
			"org":     r.GetOrg(),
			"repo":    r.GetName(),
			"repo_id": r.GetID(),
		}).Info("hook updated")
	}()

	// attach a sender SCM id if the webhook payload from the SCM has no sender id
	// the code in ProcessWebhook implies that the sender may not always be present
	// fallbacks like pusher/commit_author do not have an id
	if len(b.GetSenderSCMID()) == 0 || b.GetSenderSCMID() == "0" {
		// fetch scm user id for pusher
		senderID, err := scm.FromContext(c).GetUserID(ctx, b.GetSender(), repo.GetOwner().GetToken())
		if err != nil {
			retErr := fmt.Errorf("unable to assign sender SCM id: %w", err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}

		b.SetSenderSCMID(senderID)
	}

	// set the RepoID fields
	b.SetRepo(repo)
	h.SetRepo(repo)

	// number of times to retry
	retryLimit := 3
	// implement a loop to process asynchronous operations with a retry limit
	//
	// Some operations taken during the webhook workflow can lead to race conditions
	// failing to successfully process the request. This logic ensures we attempt our
	// best efforts to handle these cases gracefully.
	for i := 0; i < retryLimit; i++ {
		// check if we're on the first iteration of the loop
		if i > 0 {
			// incrementally sleep in between retries
			time.Sleep(time.Duration(i) * time.Second)
		}

		// send API call to capture the last hook for the repo
		lastHook, err := database.FromContext(c).LastHookForRepo(ctx, repo)
		if err != nil {
			// format the error message with extra information
			err = fmt.Errorf("unable to get last hook for repo %s: %w", r.GetFullName(), err)

			// log the error for traceability
			logrus.Error(err.Error())

			// check if the retry limit has been exceeded
			if i < retryLimit {
				// continue to the next iteration of the loop
				continue
			}

			retErr := fmt.Errorf("%s: %w", baseErr, err)
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
			// format the error message with extra information
			err = fmt.Errorf("unable to create webhook %s/%d: %w", r.GetFullName(), h.GetNumber(), err)

			// log the error for traceability
			logrus.Error(err.Error())

			// check if the retry limit has been exceeded
			if i < retryLimit {
				// continue to the next iteration of the loop
				continue
			}

			retErr := fmt.Errorf("%s: %w", baseErr, err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}

		// hook was created successfully
		break
	}

	l.WithFields(logrus.Fields{
		"hook":    h.GetNumber(),
		"hook_id": h.GetID(),
		"org":     repo.GetOrg(),
		"repo":    repo.GetName(),
	}).Info("hook created")

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
		Build:      b,
		Deployment: webhook.Deployment,
		Metadata:   m,
		BaseErr:    baseErr,
		Source:     "webhook",
		Comment:    prComment,
		Labels:     prLabels,
		Retries:    3,
	}

	// generate the queue item
	p, item, code, err := build.CompileAndPublish(
		c,
		config,
		database.FromContext(c),
		scm.FromContext(c),
		compiler.FromContext(c),
		queue.FromContext(c),
	)

	// error handling done in CompileAndPublish
	if err != nil && code == http.StatusOK {
		h.SetStatus(constants.StatusSkipped)
		h.SetError(err.Error())

		c.JSON(http.StatusOK, err.Error())

		return
	}

	if err != nil {
		h.SetStatus(constants.StatusFailure)
		h.SetError(err.Error())

		b.SetStatus(constants.StatusError)

		util.HandleError(c, code, err)

		err = scm.FromContext(c).Status(ctx, repo.GetOwner(), b, repo.GetOrg(), repo.GetName())
		if err != nil {
			l.Debugf("unable to set commit status for %s/%d: %v", repo.GetFullName(), b.GetNumber(), err)
		}

		return
	}

	// capture the build and repo from the items
	b = item.Build

	// set hook build
	h.SetBuild(b)

	// if event is deployment, update the deployment record to include this build
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		d, err := database.FromContext(c).GetDeploymentForRepo(c, repo, webhook.Deployment.GetNumber())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				deployment := webhook.Deployment

				deployment.SetRepo(repo)
				deployment.SetBuilds([]*types.Build{b})

				dr, err := database.FromContext(c).CreateDeployment(c, deployment)
				if err != nil {
					retErr := fmt.Errorf("%s: failed to create deployment %s/%d: %w", baseErr, repo.GetFullName(), deployment.GetNumber(), err)
					util.HandleError(c, http.StatusInternalServerError, retErr)

					h.SetStatus(constants.StatusFailure)
					h.SetError(retErr.Error())

					return
				}

				l.WithFields(logrus.Fields{
					"deployment":    dr.GetNumber(),
					"deployment_id": dr.GetID(),
					"org":           repo.GetOrg(),
					"repo":          repo.GetName(),
					"repo_id":       repo.GetID(),
				}).Info("deployment created")
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

			l.WithFields(logrus.Fields{
				"deployment":    d.GetNumber(),
				"deployment_id": d.GetID(),
				"org":           repo.GetOrg(),
				"repo":          repo.GetName(),
				"repo_id":       repo.GetID(),
			}).Info("deployment updated")
		}
	}

	// regardless of whether the build is published to queue, we want to attempt to auto-cancel if no errors
	defer func() {
		if err == nil && build.ShouldAutoCancel(p.Metadata.AutoCancel, b, repo.GetBranch()) {
			// fetch pending and running builds
			rBs, err := database.FromContext(c).ListPendingAndRunningBuildsForRepo(c, repo)
			if err != nil {
				l.Errorf("unable to fetch pending and running builds for %s: %v", repo.GetFullName(), err)
			}

			l.WithFields(logrus.Fields{
				"build":    b.GetNumber(),
				"build_id": b.GetID(),
				"org":      repo.GetOrg(),
				"repo":     repo.GetName(),
				"repo_id":  repo.GetID(),
			}).Debugf("found %d pending/running builds", len(rBs))

			for _, rB := range rBs {
				// call auto cancel routine
				canceled, err := build.AutoCancel(c, b, rB, p.Metadata.AutoCancel)
				if err != nil {
					// continue cancel loop if error, but log based on type of error
					if canceled {
						l.Errorf("unable to update canceled build error message: %v", err)
					} else {
						l.Errorf("unable to cancel running build: %v", err)
					}
				}

				if canceled {
					l.WithFields(logrus.Fields{
						"build":    rB.GetNumber(),
						"build_id": rB.GetID(),
						"org":      repo.GetOrg(),
						"repo":     repo.GetName(),
						"repo_id":  repo.GetID(),
					}).Debug("auto-canceled build")
				}
			}
		}
	}()

	// determine whether to send compiled build to queue
	shouldEnqueue, err := build.ShouldEnqueue(c, l, b, repo)
	if err != nil {
		retErr := fmt.Errorf("unable to process build destination: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return
	}

	if shouldEnqueue {
		// send API call to set the status on the commit
		err := scm.FromContext(c).Status(c.Request.Context(), repo.GetOwner(), b, repo.GetOrg(), repo.GetName())
		if err != nil {
			l.Errorf("unable to set commit status for %s/%d: %v", repo.GetFullName(), b.GetNumber(), err)
		}

		// publish the build to the queue
		go build.Enqueue(
			context.WithoutCancel(c.Request.Context()),
			queue.FromGinContext(c),
			database.FromContext(c),
			item,
			b.GetHost(),
		)
	} else {
		err := build.GatekeepBuild(c, b, repo)
		if err != nil {
			retErr := fmt.Errorf("unable to gate build: %w", err)
			util.HandleError(c, http.StatusInternalServerError, err)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return
		}
	}

	c.JSON(http.StatusCreated, b)
}

// handleRepositoryEvent is a helper function that processes repository events from the SCM and updates
// the database resources with any relevant changes resulting from the event, such as name changes, transfers, etc.
//
// the caller is responsible for returning errors to the client.
func handleRepositoryEvent(ctx context.Context, l *logrus.Entry, db database.Interface, m *internal.Metadata, h *types.Hook, r *types.Repo, dbRepo *types.Repo) (*types.Repo, error) {
	l = l.WithFields(logrus.Fields{
		"event_type": h.GetEvent(),
	})

	l.Debugf("webhook is repository event, making necessary updates to repo %s", r.GetFullName())

	defer func() {
		// send API call to update the webhook
		hr, err := db.CreateHook(ctx, h)
		if err != nil {
			l.Errorf("unable to create webhook %s/%d: %v", r.GetFullName(), h.GetNumber(), err)
		}

		l.WithFields(logrus.Fields{
			"hook":    hr.GetNumber(),
			"hook_id": hr.GetID(),
			"org":     r.GetOrg(),
			"repo":    r.GetName(),
			"repo_id": r.GetID(),
		}).Info("hook created")
	}()

	switch h.GetEventAction() {
	// if action is renamed or transferred, go through rename routine
	case constants.ActionRenamed, constants.ActionTransferred:
		r, err := RenameRepository(ctx, l, db, h, r, dbRepo, m)
		if err != nil {
			h.SetStatus(constants.StatusFailure)
			h.SetError(err.Error())

			return nil, err
		}

		return r, nil
	// if action is archived, unarchived, or edited, perform edits to relevant repo fields
	case "archived", "unarchived", constants.ActionEdited:
		l.Debugf("repository action %s for %s", h.GetEventAction(), r.GetFullName())

		// send API call to capture the last hook for the repo
		lastHook, err := db.LastHookForRepo(ctx, dbRepo)
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

		h.SetRepo(dbRepo)

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
		dbRepo, err = db.UpdateRepo(ctx, dbRepo)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to update repo %s: %w", baseErr, r.GetFullName(), err)

			h.SetStatus(constants.StatusFailure)
			h.SetError(retErr.Error())

			return nil, err
		}

		l.WithFields(logrus.Fields{
			"org":     dbRepo.GetOrg(),
			"repo":    dbRepo.GetName(),
			"repo_id": dbRepo.GetID(),
		}).Info("repo updated")

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
//
// the caller is responsible for returning errors to the client.
func RenameRepository(ctx context.Context, l *logrus.Entry, db database.Interface, h *types.Hook, r *types.Repo, dbR *types.Repo, m *internal.Metadata) (*types.Repo, error) {
	l = l.WithFields(logrus.Fields{
		"event_type": h.GetEvent(),
	})

	l.Debugf("renaming repository from %s to %s", r.GetPreviousName(), r.GetName())

	// update hook object which will be added to DB upon reaching deferred function in PostWebhook
	h.SetRepo(r)

	// send API call to capture the last hook for the repo
	lastHook, err := db.LastHookForRepo(ctx, dbR)
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

	// get total number of secrets associated with repository
	t, err := db.CountSecretsForRepo(ctx, dbR, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unable to get secret count for repo %s/%s: %w", dbR.GetOrg(), dbR.GetName(), err)
	}

	secrets := []*types.Secret{}
	page := 1
	// capture all secrets belonging to certain repo in database
	for repoSecrets := int64(0); repoSecrets < t; repoSecrets += 100 {
		s, err := db.ListSecretsForRepo(ctx, dbR, map[string]interface{}{}, page, 100)
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

		_, err = db.UpdateSecret(ctx, secret)
		if err != nil {
			return nil, fmt.Errorf("unable to update secret for repo %s/%s: %w", dbR.GetOrg(), dbR.GetName(), err)
		}

		l.WithFields(logrus.Fields{
			"secret_id": secret.GetID(),
			"repo":      secret.GetRepo(),
			"org":       secret.GetOrg(),
		}).Info("secret updated")
	}

	// get total number of builds associated with repository
	t, err = db.CountBuildsForRepo(ctx, dbR, nil, time.Now().Unix(), 0)
	if err != nil {
		return nil, fmt.Errorf("unable to get build count for repo %s: %w", dbR.GetFullName(), err)
	}

	builds := []*types.Build{}
	page = 1
	// capture all builds belonging to repo in database
	for build := int64(0); build < t; build += 100 {
		b, err := db.ListBuildsForRepo(ctx, dbR, nil, time.Now().Unix(), 0, page, 100)
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

		_, err = db.UpdateBuild(ctx, build)
		if err != nil {
			return nil, fmt.Errorf("unable to update build for repo %s: %w", dbR.GetFullName(), err)
		}

		l.WithFields(logrus.Fields{
			"build_id": build.GetID(),
			"build":    build.GetNumber(),
			"org":      dbR.GetOrg(),
			"repo":     dbR.GetName(),
			"repo_id":  dbR.GetID(),
		}).Info("build updated")
	}

	// update the repo name information
	dbR.SetName(r.GetName())
	dbR.SetOrg(r.GetOrg())
	dbR.SetFullName(r.GetFullName())
	dbR.SetClone(r.GetClone())
	dbR.SetLink(r.GetLink())
	dbR.SetPreviousName(r.GetPreviousName())

	// update the repo in the database
	dbR, err = db.UpdateRepo(ctx, dbR)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to update repo %s/%s", baseErr, dbR.GetOrg(), dbR.GetName())

		h.SetStatus(constants.StatusFailure)
		h.SetError(retErr.Error())

		return nil, retErr
	}

	l.WithFields(logrus.Fields{
		"org":     dbR.GetOrg(),
		"repo":    dbR.GetName(),
		"repo_id": dbR.GetID(),
	}).Infof("repo updated in database (previous name: %s)", r.GetPreviousName())

	return dbR, nil
}
