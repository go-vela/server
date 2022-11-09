// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
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
	if len(input.GetVisibility()) == 0 {
		// default visibility field to public
		r.SetVisibility(constants.VisibilityPublic)
	} else {
		r.SetVisibility(input.GetVisibility())
	}

	// set default events if no events are passed in
	if !input.GetAllowPull() && !input.GetAllowPush() &&
		!input.GetAllowDeploy() && !input.GetAllowTag() &&
		!input.GetAllowComment() {
		// default event to push
		r.SetAllowPush(true)
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
	if !checkAllowlist(r, allowlist) {
		retErr := fmt.Errorf("unable to activate repo: %s is not on allowlist", r.GetFullName())

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	// send API call to capture the repo from the database
	dbRepo, err := database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())
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

	// send API call to create the webhook
	if c.Value("webhookvalidation").(bool) {
		_, err = scm.FromContext(c).Enable(u, r.GetOrg(), r.GetName(), r.GetHash())
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
		err = database.FromContext(c).UpdateRepo(dbRepo)
		if err != nil {
			retErr := fmt.Errorf("unable to set repo %s to active: %w", dbRepo.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture the updated repo
		r, _ = database.FromContext(c).GetRepoForOrg(dbRepo.GetOrg(), dbRepo.GetName())
	} else {
		// send API call to create the repo
		err = database.FromContext(c).CreateRepo(r)
		if err != nil {
			retErr := fmt.Errorf("unable to create new repo %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture the created repo
		r, _ = database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())
	}

	c.JSON(http.StatusCreated, r)
}

// swagger:operation GET /api/v1/repos repos GetRepos
//
// Get all repos in the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
// - in: query
//   name: page
//   description: The page of results to retrieve
//   type: integer
//   default: 1
// - in: query
//   name: per_page
//   description: How many results per page to return
//   type: integer
//   maximum: 100
//   default: 10
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Repo"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the repo
//     schema:
//       "$ref": "#/definitions/Error"

// GetRepos represents the API handler to capture a list
// of repos for a user from the configured backend.
func GetRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("reading repos for user %s", u.GetName())

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// capture the sort_by query parameter if present
	sortBy := util.QueryParameter(c, "sort_by", "name")

	// capture the query parameters if present:
	//
	// * active
	filters := map[string]interface{}{
		"active": util.QueryParameter(c, "active", "true"),
	}

	// send API call to capture the list of repos for the user
	r, t, err := database.FromContext(c).ListReposForUser(u, sortBy, filters, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get repos for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, r)
}

// swagger:operation GET /api/v1/repos/{org} repos GetOrgRepos
//
// Get all repos for the provided org in the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: query
//   name: active
//   description: Filter active repos
//   type: boolean
//   default: true
// - in: query
//   name: page
//   description: The page of results to retrieve
//   type: integer
//   default: 1
// - in: query
//   name: per_page
//   description: How many results per page to return
//   type: integer
//   maximum: 100
//   default: 10
// - in: query
//   name: sort_by
//   description: How to sort the results
//   type: string
//   enum:
//   - name
//   - latest
//   default: name
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Repo"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the org
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the org
//     schema:
//       "$ref": "#/definitions/Error"

// GetOrgRepos represents the API handler to capture a list
// of repos for an org from the configured backend.
func GetOrgRepos(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"user": u.GetName(),
	}).Infof("reading repos for org %s", o)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// capture the sort_by query parameter if present
	sortBy := util.QueryParameter(c, "sort_by", "name")

	// capture the query parameters if present:
	//
	// * active
	filters := map[string]interface{}{
		"active": util.QueryParameter(c, "active", "true"),
	}

	// See if the user is an org admin to bypass individual permission checks
	perm, err := scm.FromContext(c).OrgAccess(u, o)
	if err != nil {
		logrus.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}
	// Only show public repos to non-admins
	if perm != "admin" {
		filters["visibility"] = constants.VisibilityPublic
	}

	// send API call to capture the list of repos for the org
	r, t, err := database.FromContext(c).ListReposForOrg(o, sortBy, filters, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get repos for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, r)
}

// swagger:operation GET /api/v1/repos/{org}/{repo} repos GetRepo
//
// Get a repo in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     schema:
//       "$ref": "#/definitions/Repo"

// GetRepo represents the API handler to
// capture a repo from the configured backend.
func GetRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("reading repo %s", r.GetFullName())

	c.JSON(http.StatusOK, r)
}

// swagger:operation PUT /api/v1/repos/{org}/{repo} repos UpdateRepo
//
// Update a repo in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the repo to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the repo
//     schema:
//       "$ref": "#/definitions/Repo"
//   '400':
//     description: Unable to update the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to update the repo
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateRepo represents the API handler to update
// a repo in the configured backend.
//
//nolint:funlen // ignore line length
func UpdateRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	maxBuildLimit := c.Value("maxBuildLimit").(int64)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("updating repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update repo fields if provided
	if len(input.GetBranch()) > 0 {
		// update branch if set
		r.SetBranch(input.GetBranch())
	}

	// update build limit if set
	if input.GetBuildLimit() > 0 {
		// allow build limit between 1 - value configured by server
		r.SetBuildLimit(
			int64(
				util.MaxInt(
					constants.BuildLimitMin,
					util.MinInt(
						int(input.GetBuildLimit()),
						int(maxBuildLimit),
					), // clamp max
				), // clamp min
			),
		)
	}

	if input.GetTimeout() > 0 {
		// update build timeout if set
		r.SetTimeout(
			int64(
				util.MaxInt(
					constants.BuildTimeoutMin,
					util.MinInt(
						int(input.GetTimeout()),
						constants.BuildTimeoutMax,
					), // clamp max
				), // clamp min
			),
		)
	}

	if input.GetCounter() > 0 {
		if input.GetCounter() <= r.GetCounter() {
			retErr := fmt.Errorf("unable to set counter for repo %s: must be greater than current %d",
				r.GetFullName(), r.GetCounter())

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		r.SetCounter(input.GetCounter())
	}

	if len(input.GetVisibility()) > 0 {
		// update visibility if set
		r.SetVisibility(input.GetVisibility())
	}

	if input.Private != nil {
		// update private if set
		r.SetPrivate(input.GetPrivate())
	}

	if input.Active != nil {
		// update active if set
		r.SetActive(input.GetActive())
	}

	if input.AllowPull != nil {
		// update allow_pull if set
		r.SetAllowPull(input.GetAllowPull())
	}

	if input.AllowPush != nil {
		// update allow_push if set
		r.SetAllowPush(input.GetAllowPush())
	}

	if input.AllowDeploy != nil {
		// update allow_deploy if set
		r.SetAllowDeploy(input.GetAllowDeploy())
	}

	if input.AllowTag != nil {
		// update allow_tag if set
		r.SetAllowTag(input.GetAllowTag())
	}

	if input.AllowComment != nil {
		// update allow_comment if set
		r.SetAllowComment(input.GetAllowComment())
	}

	// set default events if no events are enabled
	if !r.GetAllowPull() && !r.GetAllowPush() &&
		!r.GetAllowDeploy() && !r.GetAllowTag() &&
		!r.GetAllowComment() {
		r.SetAllowPull(true)
		r.SetAllowPush(true)
	}

	if len(input.GetPipelineType()) != 0 {
		// ensure the pipeline type matches one of the expected values
		if input.GetPipelineType() != constants.PipelineTypeYAML &&
			input.GetPipelineType() != constants.PipelineTypeGo &&
			input.GetPipelineType() != constants.PipelineTypeStarlark {
			retErr := fmt.Errorf("pipeline_type of %s is invalid", input.GetPipelineType())

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		r.SetPipelineType(input.GetPipelineType())
	}

	// set hash for repo if no hash is already set
	if len(r.GetHash()) == 0 {
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
	}

	// send API call to update the repo
	err = database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to update repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated repo
	r, _ = database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())

	c.JSON(http.StatusOK, r)
}

// swagger:operation DELETE /api/v1/repos/{org}/{repo} repos DeleteRepo
//
// Delete a repo in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to  deleted the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '510':
//     description: Unable to  deleted the repo
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteRepo represents the API handler to remove
// a repo from the configured backend.
func DeleteRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("deleting repo %s", r.GetFullName())

	// send API call to remove the webhook
	err := scm.FromContext(c).Disable(u, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to delete webhook for %s: %w", r.GetFullName(), err)

		if err.Error() == "Repo not found" {
			util.HandleError(c, http.StatusNotExtended, retErr)

			return
		}

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Mark the the repo as inactive
	r.SetActive(false)

	err = database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to set repo %s to inactive: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Comment out actual delete until delete mechanism is fleshed out
	// err = database.FromContext(c).DeleteRepo(r.ID)
	// if err != nil {
	// 	retErr := fmt.Errorf("Error while deleting repo %s: %w", r.FullName, err)
	// 	util.HandleError(c, http.StatusInternalServerError, retErr)
	// 	return
	// }

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s deleted", r.GetFullName()))
}

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/repair repos RepairRepo
//
// Remove and recreate the webhook for a repo
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully repaired the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to repair the repo
//     schema:
//       "$ref": "#/definitions/Error"

// RepairRepo represents the API handler to remove
// and then create a webhook for a repo.
func RepairRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("repairing repo %s", r.GetFullName())

	// send API call to remove the webhook
	err := scm.FromContext(c).Disable(u, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to delete webhook for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to create the webhook
	_, err = scm.FromContext(c).Enable(u, r.GetOrg(), r.GetName(), r.GetHash())
	if err != nil {
		retErr := fmt.Errorf("unable to create webhook for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// if the repo was previously inactive, mark it as active
	if !r.GetActive() {
		r.SetActive(true)

		// send API call to update the repo
		err = database.FromContext(c).UpdateRepo(r)
		if err != nil {
			retErr := fmt.Errorf("unable to set repo %s to active: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s repaired", r.GetFullName()))
}

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/chown repos ChownRepo
//
// Change the owner of the webhook for a repo
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully changed the owner for the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to change the owner for the repo
//     schema:
//       "$ref": "#/definitions/Error"

// ChownRepo represents the API handler to change
// the owner of a repo in the configured backend.
func ChownRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("changing owner of repo %s to %s", r.GetFullName(), u.GetName())

	// update repo owner
	r.SetUserID(u.GetID())

	// send API call to updated the repo
	err := database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to change owner of repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s changed owner", r.GetFullName()))
}

// checkAllowlist is a helper function to ensure only repos in the
// allowlist are allowed to enable repos.
//
// a single entry of '*' allows any repo to be enabled.
func checkAllowlist(r *library.Repo, allowlist []string) bool {
	// check if all repos are allowed to be enabled
	if len(allowlist) == 1 && allowlist[0] == "*" {
		return true
	}

	for _, repo := range allowlist {
		// allow all repos in org
		if strings.Contains(repo, "/*") {
			if strings.HasPrefix(repo, r.GetOrg()) {
				return true
			}
		}

		// allow specific repo within org
		if repo == r.GetFullName() {
			return true
		}
	}

	return false
}
