// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/settings"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
)

// swagger:operation POST /api/v1/repos repos CreateRepo
//
// Create a repository
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Repo object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the repo
//     schema:
//       "$ref": "#/definitions/Repo"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '403':
//     description: Unable to create the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '409':
//     description: Unable to create the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to create the repo
//     schema:
//       "$ref": "#/definitions/Error"

// CreateRepo represents the API handler to create a repository.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func CreateRepo(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)
	s := settings.FromContext(c)

	defaultBuildLimit := c.Value("defaultBuildLimit").(int64)
	defaultTimeout := c.Value("defaultTimeout").(int64)
	maxBuildLimit := c.Value("maxBuildLimit").(int64)
	defaultRepoEvents := c.Value("defaultRepoEvents").([]string)
	defaultRepoEventsMask := c.Value("defaultRepoEventsMask").(int64)
	defaultRepoApproveBuild := c.Value("defaultRepoApproveBuild").(string)

	ctx := c.Request.Context()

	// capture body from API request
	input := new(types.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new repo: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	l.Debugf("creating new repo %s", input.GetFullName())

	// get repo information from the source
	r, _, err := scm.FromContext(c).GetRepo(ctx, u, input)
	if err != nil {
		retErr := fmt.Errorf("unable to retrieve repo info for %s from source: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in repo object
	r.SetOwner(u)

	// set the active field based off the input provided
	if input.Active == nil {
		// default active field to true
		r.SetActive(true)
	} else {
		r.SetActive(input.GetActive())
	}

	// set the build limit field based off the input provided
	if input.GetBuildLimit() == 0 {
		// default build limit to value configured by server
		r.SetBuildLimit(defaultBuildLimit)
	} else if input.GetBuildLimit() > maxBuildLimit {
		// set build limit to value configured by server to prevent limit from exceeding max
		r.SetBuildLimit(maxBuildLimit)
	} else {
		r.SetBuildLimit(input.GetBuildLimit())
	}

	// set the timeout field based off the input provided
	if input.GetTimeout() == 0 && defaultTimeout == 0 {
		// default build timeout to 30m
		r.SetTimeout(constants.BuildTimeoutDefault)
	} else if input.GetTimeout() == 0 {
		r.SetTimeout(defaultTimeout)
	} else {
		r.SetTimeout(input.GetTimeout())
	}

	// set the visibility field based off the input provided
	if len(input.GetVisibility()) > 0 {
		// default visibility field to the input visibility
		r.SetVisibility(input.GetVisibility())
	}

	// set the fork policy field based off the input provided
	if len(input.GetApproveBuild()) > 0 {
		// ensure the approve build setting matches one of the expected values
		if input.GetApproveBuild() != constants.ApproveForkAlways &&
			input.GetApproveBuild() != constants.ApproveForkNoWrite &&
			input.GetApproveBuild() != constants.ApproveNever &&
			input.GetApproveBuild() != constants.ApproveOnce {
			retErr := fmt.Errorf("approve_build of %s is invalid", input.GetApproveBuild())

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		r.SetApproveBuild(input.GetApproveBuild())
	} else {
		r.SetApproveBuild(defaultRepoApproveBuild)
	}

	// fields restricted to platform admins
	if u.GetAdmin() {
		// trusted default is false
		if input.GetTrusted() != r.GetTrusted() {
			r.SetTrusted(input.GetTrusted())
		}
	}

	// set allow events based on input if given
	if input.GetAllowEvents().ToDatabase() != 0 {
		r.SetAllowEvents(input.GetAllowEvents())
	} else {
		r.SetAllowEvents(defaultAllowedEvents(defaultRepoEvents, defaultRepoEventsMask))
	}

	if len(input.GetPipelineType()) == 0 {
		r.SetPipelineType(constants.PipelineTypeYAML)
	} else {
		// ensure the pipeline type matches one of the expected values
		if input.GetPipelineType() != constants.PipelineTypeYAML &&
			input.GetPipelineType() != constants.PipelineTypeGo &&
			input.GetPipelineType() != constants.PipelineTypeStarlark {
			retErr := fmt.Errorf("unable to create new repo %s: invalid pipeline_type provided %s", r.GetFullName(), input.GetPipelineType())

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		r.SetPipelineType(input.GetPipelineType())
	}

	// create unique id for the repo
	uid, err := uuid.NewRandom()
	if err != nil {
		retErr := fmt.Errorf("unable to create UID for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	r.SetHash(
		base64.StdEncoding.EncodeToString(
			[]byte(strings.TrimSpace(uid.String())),
		),
	)

	// ensure repo is allowed to be activated
	if !util.CheckAllowlist(r, s.GetRepoAllowlist()) {
		retErr := fmt.Errorf("unable to activate repo: %s is not on allowlist", r.GetFullName())

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	// send API call to capture the repo from the database
	dbRepo, err := database.FromContext(c).GetRepoForOrg(ctx, r.GetFullName())
	if err == nil && dbRepo.GetActive() {
		retErr := fmt.Errorf("unable to activate repo: %s is already active", r.GetFullName())

		util.HandleError(c, http.StatusConflict, retErr)

		return
	}

	// check if the repo already has a hash created
	if len(dbRepo.GetHash()) > 0 {
		// overwrite the new repo hash with the existing repo hash
		r.SetHash(dbRepo.GetHash())
	}

	h := new(types.Hook)

	// err being nil means we have a record of this repo (dbRepo)
	if err == nil {
		h, _ = database.FromContext(c).LastHookForRepo(ctx, dbRepo)

		// make sure our record of the repo allowed events matches what we send to SCM
		// what the dbRepo has should override default events on enable
		r.SetAllowEvents(dbRepo.GetAllowEvents())
	}

	// check if we should create the webhook
	if c.Value("webhookvalidation").(bool) {
		// send API call to create the webhook
		h, _, err = scm.FromContext(c).Enable(ctx, u, r, h)
		if err != nil {
			retErr := fmt.Errorf("unable to create webhook for %s: %w", r.GetFullName(), err)

			switch err.Error() {
			case "repo already enabled":
				util.HandleError(c, http.StatusConflict, retErr)
				return
			case "repo not found":
				util.HandleError(c, http.StatusNotFound, retErr)
				return
			}

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// if the repo exists but is inactive
	if len(dbRepo.GetOrg()) > 0 && !dbRepo.GetActive() {
		// update the repo owner
		dbRepo.SetOwner(u)
		// update the default branch
		dbRepo.SetBranch(r.GetBranch())
		// activate the repo
		dbRepo.SetActive(true)

		// send API call to update the repo
		// NOTE: not logging modification out separately
		// although we are CREATING a repo in this path
		r, err = database.FromContext(c).UpdateRepo(ctx, dbRepo)
		if err != nil {
			retErr := fmt.Errorf("unable to set repo %s to active: %w", dbRepo.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.WithFields(logrus.Fields{
			"org":     r.GetOrg(),
			"repo":    r.GetName(),
			"repo_id": r.GetID(),
		}).Infof("repo %s activated", r.GetFullName())
	} else {
		// send API call to create the repo
		r, err = database.FromContext(c).CreateRepo(ctx, r)
		if err != nil {
			retErr := fmt.Errorf("unable to create new repo %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.WithFields(logrus.Fields{
			"org":     r.GetOrg(),
			"repo":    r.GetName(),
			"repo_id": r.GetID(),
		}).Infof("repo %s created", r.GetFullName())
	}

	// create init hook in the DB after repo has been added in order to capture its ID
	if c.Value("webhookvalidation").(bool) {
		// update initialization hook
		h.SetRepo(r)
		// create first hook for repo in the database
		_, err = database.FromContext(c).CreateHook(ctx, h)
		if err != nil {
			retErr := fmt.Errorf("unable to create initialization webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.WithFields(logrus.Fields{
			"hook": h.GetID(),
		}).Infof("hook %d created for repo %s", h.GetID(), r.GetFullName())
	}

	c.JSON(http.StatusCreated, r)
}

// defaultAllowedEvents is a helper function that generates an Events struct that results
// from an admin-provided `sliceDefaults` or an admin-provided `maskDefaults`. If the admin
// supplies a mask, that will be the default. Otherwise, it will be the legacy event list.
func defaultAllowedEvents(sliceDefaults []string, maskDefaults int64) *types.Events {
	if maskDefaults > 0 {
		return types.NewEventsFromMask(maskDefaults)
	}

	events := new(types.Events)

	for _, event := range sliceDefaults {
		switch event {
		case constants.EventPull:
			pull := new(actions.Pull)
			pull.SetOpened(true)
			pull.SetSynchronize(true)

			events.SetPullRequest(pull)
		case constants.EventPush:
			push := events.GetPush()
			push.SetBranch(true)

			events.SetPush(push)
		case constants.EventTag:
			tag := events.GetPush()
			tag.SetTag(true)

			events.SetPush(tag)
		case constants.EventDeploy:
			deploy := new(actions.Deploy)
			deploy.SetCreated(true)

			events.SetDeployment(deploy)
		case constants.EventComment:
			comment := new(actions.Comment)
			comment.SetCreated(true)
			comment.SetEdited(true)

			events.SetComment(comment)
		case constants.EventDelete:
			deletion := events.GetPush()
			deletion.SetDeleteBranch(true)
			deletion.SetDeleteTag(true)

			events.SetPush(deletion)
		}
	}

	return events
}
