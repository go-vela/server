// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos repos CreateRepo
//
// Create a repo in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the repo to create
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
//     description: Unable to create the repo
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
//     description: Unable to create the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to create the repo
//     schema:
//       "$ref": "#/definitions/Error"

// CreateRepo represents the API handler to
// create a repo in the configured backend.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func CreateRepo(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	allowlist := c.Value("allowlist").([]string)
	defaultBuildLimit := c.Value("defaultBuildLimit").(int64)
	defaultTimeout := c.Value("defaultTimeout").(int64)
	maxBuildLimit := c.Value("maxBuildLimit").(int64)
	defaultRepoEvents := c.Value("defaultRepoEvents").([]string)
	ctx := c.Request.Context()

	// capture body from API request
	input := new(library.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new repo: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  input.GetOrg(),
		"repo": input.GetName(),
		"user": u.GetName(),
	}).Infof("creating new repo %s", input.GetFullName())

	// get repo information from the source
	r, err := scm.FromContext(c).GetRepo(u, input)
	if err != nil {
		retErr := fmt.Errorf("unable to retrieve repo info for %s from source: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in repo object
	r.SetUserID(u.GetID())

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

	// fields restricted to platform admins
	if u.GetAdmin() {
		// trusted default is false
		if input.GetTrusted() != r.GetTrusted() {
			r.SetTrusted(input.GetTrusted())
		}
	}

	// set default events if no events are passed in
	if !input.GetAllowPull() && !input.GetAllowPush() &&
		!input.GetAllowDeploy() && !input.GetAllowTag() &&
		!input.GetAllowComment() {
		for _, event := range defaultRepoEvents {
			switch event {
			case constants.EventPull:
				r.SetAllowPull(true)
			case constants.EventPush:
				r.SetAllowPush(true)
			case constants.EventDeploy:
				r.SetAllowDeploy(true)
			case constants.EventTag:
				r.SetAllowTag(true)
			case constants.EventComment:
				r.SetAllowComment(true)
			}
		}
	} else {
		r.SetAllowComment(input.GetAllowComment())
		r.SetAllowDeploy(input.GetAllowDeploy())
		r.SetAllowPull(input.GetAllowPull())
		r.SetAllowPush(input.GetAllowPush())
		r.SetAllowTag(input.GetAllowTag())
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
	if !util.CheckAllowlist(r, allowlist) {
		retErr := fmt.Errorf("unable to activate repo: %s is not on allowlist", r.GetFullName())

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	// send API call to capture the repo from the database
	dbRepo, err := database.FromContext(c).GetRepoForOrg(ctx, r.GetOrg(), r.GetName())
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

	h := new(library.Hook)

	// err being nil means we have a record of this repo (dbRepo)
	if err == nil {
		h, _ = database.FromContext(c).LastHookForRepo(ctx, dbRepo)

		// make sure our record of the repo allowed events matches what we send to SCM
		// what the dbRepo has should override default events on enable
		r.SetAllowComment(dbRepo.GetAllowComment())
		r.SetAllowDeploy(dbRepo.GetAllowDeploy())
		r.SetAllowPull(dbRepo.GetAllowPull())
		r.SetAllowPush(dbRepo.GetAllowPush())
		r.SetAllowTag(dbRepo.GetAllowTag())
	}

	// check if we should create the webhook
	if c.Value("webhookvalidation").(bool) {
		// send API call to create the webhook
		h, _, err = scm.FromContext(c).Enable(u, r, h)
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
		dbRepo.SetUserID(u.GetID())
		// update the default branch
		dbRepo.SetBranch(r.GetBranch())
		// activate the repo
		dbRepo.SetActive(true)

		// send API call to update the repo
		r, err = database.FromContext(c).UpdateRepo(ctx, dbRepo)
		if err != nil {
			retErr := fmt.Errorf("unable to set repo %s to active: %w", dbRepo.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	} else {
		// send API call to create the repo
		r, err = database.FromContext(c).CreateRepo(ctx, r)
		if err != nil {
			retErr := fmt.Errorf("unable to create new repo %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// create init hook in the DB after repo has been added in order to capture its ID
	if c.Value("webhookvalidation").(bool) {
		// update initialization hook
		h.SetRepoID(r.GetID())
		// create first hook for repo in the database
		_, err = database.FromContext(c).CreateHook(ctx, h)
		if err != nil {
			retErr := fmt.Errorf("unable to create initialization webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	c.JSON(http.StatusCreated, r)
}
