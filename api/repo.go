// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos repos CreateRepo
//
// Create a repo in the configured backend
//
// ---
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the repo to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '201':
//     description: Successfully created the repo
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"
//   '400':
//     description: Unable to create the repo
//     schema:
//       type: string
//   '403':
//     description: Unable to create the repo
//     schema:
//       type: string
//   '409':
//     description: Unable to create the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to create the repo
//     schema:
//       type: string
//   '503':
//     description: Unable to create the repo
//     schema:
//       type: string

// CreateRepo represents the API handler to
// create a repo in the configured backend.
func CreateRepo(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	whitelist := c.Value("whitelist").([]string)

	logrus.Info("Creating new repo")

	// capture body from API request
	input := new(library.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new repo: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in repo object
	input.SetUserID(u.GetID())

	if input.Active == nil {
		input.SetActive(true)
	}

	if input.GetTimeout() == 0 {
		input.SetTimeout(constants.BuildTimeoutMin)
	}

	if len(input.GetVisibility()) == 0 {
		input.SetVisibility(constants.VisibilityPublic)
	}

	if len(input.GetFullName()) == 0 {
		input.SetFullName(fmt.Sprintf("%s/%s", input.GetOrg(), input.GetName()))
	}

	if len(input.GetBranch()) == 0 {
		input.SetBranch("master")
	}

	if !input.GetAllowPull() && !input.GetAllowPush() &&
		!input.GetAllowDeploy() && !input.GetAllowTag() &&
		!input.GetAllowComment() {
		input.SetAllowPull(true)
		input.SetAllowPush(true)
	}

	// create unique id for the repo
	uid, err := uuid.NewRandom()
	if err != nil {
		retErr := fmt.Errorf("unable to create UID for repo %s: %w", input.GetFullName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	input.SetHash(
		base64.StdEncoding.EncodeToString(
			[]byte(strings.TrimSpace(uid.String())),
		),
	)

	// ensure repo is allowed to be activated
	if !checkWhitelist(input, whitelist) {
		retErr := fmt.Errorf("unable to activate repo: %s is not on whitelist", input.GetFullName())

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	// send API call to capture the repo
	r, err := database.FromContext(c).GetRepo(input.GetOrg(), input.GetName())
	if err == nil && r.GetActive() {
		retErr := fmt.Errorf("unable to activate repo: %s is already active", input.GetFullName())

		util.HandleError(c, http.StatusConflict, retErr)

		return
	}

	// check if the repo already has a hash created
	if len(r.GetHash()) > 0 {
		// overwrite the new repo hash with the existing repo hash
		input.SetHash(r.GetHash())
	}

	// send API call to create the webhook
	url, err := source.FromContext(c).Enable(u, input.GetOrg(), input.GetName(), input.GetHash())
	if err != nil {
		retErr := fmt.Errorf("unable to create webhook for %s: %w", r.GetFullName(), err)

		switch err.Error() {
		case "Repo already enabled":
			util.HandleError(c, http.StatusConflict, retErr)
			return
		case "Repo not found":
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// TODO: build these from the source client
	if len(input.GetLink()) == 0 {
		input.SetLink(url)
	}

	if len(input.GetClone()) == 0 {
		input.SetClone(fmt.Sprintf("%s.git", url))
	}

	if len(r.GetOrg()) > 0 && !r.GetActive() {
		r.SetActive(true)

		// send API call to update the repo
		err = database.FromContext(c).UpdateRepo(r)
		if err != nil {
			retErr := fmt.Errorf("unable to set repo %s to active: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture the updated repo
		r, _ = database.FromContext(c).GetRepo(r.GetOrg(), r.GetName())
	} else {
		// send API call to create the repo
		err = database.FromContext(c).CreateRepo(input)
		if err != nil {
			retErr := fmt.Errorf("unable to create new repo %s: %w", input.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture the created repo
		r, _ = database.FromContext(c).GetRepo(input.GetOrg(), input.GetName())
	}

	c.JSON(http.StatusCreated, r)
}

// swagger:operation GET /api/v1/repos repos GetRepos
//
// Get all repos in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"
//   '400':
//     description: Unable to retrieve the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve the repo
//     schema:
//       type: string

// GetRepos represents the API handler to capture a list
// of repos for a user from the configured backend.
func GetRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	logrus.Infof("Reading repos for user %s", u.GetName())

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

	// send API call to capture the total number of repos for the user
	t, err := database.FromContext(c).GetUserRepoCount(u)
	if err != nil {
		retErr := fmt.Errorf("unable to get repo count for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the list of repos for the user
	r, err := database.FromContext(c).GetUserRepoList(u, page, perPage)
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

// swagger:operation GET /api/v1/repos/{org}/{repo} repos GetRepo
//
// Get a repo in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"

// GetRepo represents the API handler to
// capture a repo from the configured backend.
func GetRepo(c *gin.Context) {
	logrus.Infof("Reading repo %s/%s", c.Param("org"), c.Param("repo"))

	// retrieve repo from context
	r := repo.Retrieve(c)

	c.JSON(http.StatusOK, r)
}

// swagger:operation PUT /api/v1/repos/{org}/{repo} repos UpdateRepo
//
// Update a repo in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the repo to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully updated the repo
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"
//   '400':
//     description: Unable to update the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to update the repo
//     schema:
//       type: string
//   '503':
//     description: Unable to update the repo
//     schema:
//       type: string
//   '510':
//     description: Unable to update the repo
//     schema:
//       type: string

// UpdateRepo represents the API handler to update
// a repo in the configured backend.
func UpdateRepo(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)

	logrus.Infof("Updating repo %s", r.GetFullName())

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
	r, _ = database.FromContext(c).GetRepo(r.GetOrg(), r.GetName())

	c.JSON(http.StatusOK, r)
}

// swagger:operation DELETE /api/v1/repos/{org}/{repo} repos DeleteRepo
//
// Delete a repo in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully deleted the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to  deleted the repo
//     schema:
//       type: string
//   '510':
//     description: Unable to  deleted the repo
//     schema:
//       type: string

// DeleteRepo represents the API handler to remove
// a repo from the configured backend.
func DeleteRepo(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	logrus.Infof("Deleting repo %s", r.GetFullName())

	// send API call to remove the webhook
	err := source.FromContext(c).Disable(u, r.GetOrg(), r.GetName())
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
	// 	retErr := fmt.Errorf("Error while deleting repo %s: %v", r.FullName, err)
	// 	util.HandleError(c, http.StatusInternalServerError, retErr)
	// 	return
	// }

	c.JSON(http.StatusOK, fmt.Sprintf("Repo %s deleted", r.GetFullName()))
}

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/repair repos RepairRepo
//
// Remove and recreate the webhook for a repo
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully repaired the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to repair the repo
//     schema:
//       type: string

// RepairRepo represents the API handler to remove
// and then create a webhook for a repo.
func RepairRepo(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	s := source.FromContext(c)

	logrus.Infof("Repairing repo %s", r.GetFullName())

	// send API call to remove the webhook
	err := s.Disable(u, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to delete webhook for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to create the webhook
	_, err = s.Enable(u, r.GetOrg(), r.GetName(), r.GetHash())
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

	c.JSON(http.StatusOK, fmt.Sprintf("Repo %s repaired", r.GetFullName()))
}

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/chown repos ChownRepo
//
// Change the owner of the webhook for a repo
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully changed the owner for the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to change the owner for the repo
//     schema:
//       type: string

// ChownRepo represents the API handler to change
// the owner of a repo in the configured backend.
func ChownRepo(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	logrus.Infof("Changing owner of repo %s", r.GetFullName())

	// update repo owner
	r.SetUserID(u.GetID())

	// send API call to updated the repo
	err := database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to change owner of repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Repo %s changed owner", r.GetFullName()))
}

// checkWhitelist is a helper function to ensure only repos in the
// whitelist are allowed to enable repos. If the whitelist is
// empty then any repo can be enabled.
func checkWhitelist(r *library.Repo, whitelist []string) bool {
	// if the whitelist is not set or empty allow any repo to be enabled
	if len(whitelist) == 0 {
		return true
	}

	for _, repo := range whitelist {
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
