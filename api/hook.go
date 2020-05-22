// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/hooks/{org}/{repo} webhook CreateHook
//
// Create a webhook for the configured backend
//
// ---
// x-success_http_code: '201'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Webhook payload that we expect from the user or VCS
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
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
//   '201':
//     description: The webhook has been created
//     schema:
//       type: string
//   '400':
//     description: The webhook was unable to be created
//     schema:
//       type: string
//   '500':
//     description: The webhook was unable to be created
//     schema:
//       type: string

// CreateHook represents the API handler to create
// a webhook in the configured backend.
func CreateHook(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)

	logrus.Infof("Creating new webhook for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new webhook for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the last hook for the repo
	lastHook, err := database.FromContext(c).GetLastHook(r)
	if err != nil {
		retErr := fmt.Errorf("unable to get last hook for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update fields in webhook object
	input.SetRepoID(r.GetID())
	input.SetNumber(1)

	if input.GetCreated() == 0 {
		input.SetCreated(time.Now().UTC().Unix())
	}

	if lastHook != nil {
		input.SetNumber(
			lastHook.GetNumber() + 1,
		)
	}

	// send API call to create the webhook
	err = database.FromContext(c).CreateHook(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create webhook for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created webhook
	h, _ := database.FromContext(c).GetHook(input.GetNumber(), r)

	c.JSON(http.StatusCreated, h)
}

// swagger:operation GET /api/v1/hooks/{org}/{repo} deployment GetHooks
//
// Create a webhook for the configured backend
//
// ---
// x-success_http_code: '200'
// x-incident_priority: P4
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
//   '201':
//     description: Successfully retrieved webhooks
//     schema:
//       type: string
//   '400':
//     description: Unable to retrieve webhooks
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve webhooks
//     schema:
//       type: string

// GetHooks represents the API handler to capture a list
// of webhooks from the configured backend.
func GetHooks(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)

	logrus.Infof("Reading webhooks for repo %s", r.GetFullName())

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the total number of webhooks for the repo
	t, err := database.FromContext(c).GetRepoHookCount(r)
	if err != nil {
		retErr := fmt.Errorf("unable to get hooks count for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the list of steps for the build
	h, err := database.FromContext(c).GetRepoHookList(r, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get hooks for repo %s: %w", r.GetFullName(), err)

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

	c.JSON(http.StatusOK, h)
}

// swagger:operation GET /api/v1/hooks/{org}/{repo}/{hook} deployment GetHook
//
// Create a webhook for the configured backend
//
// ---
// x-success_http_code: '200'
// x-incident_priority: P4
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
// - in: path
//   name: hook
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
//     description: Successfully retrieved the webhook
//     schema:
//       type: string
//   '400':
//     description: Unable to retrieve the webhook
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve the webhook
//     schema:
//       type: string

// GetHook represents the API handler to capture a
// webhook from the configured backend.
func GetHook(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	hook := c.Param("hook")

	logrus.Infof("Reading webhook %s/%s", r.GetFullName(), hook)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHook(number, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get webhook %s/%d: %w", r.GetFullName(), h.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, h)
}

// swagger:operation PUT /api/v1/hooks/{org}/{repo}/{hook} webhook UpdateHook
//
// Update a webhook for the configured backend
//
// ---
// x-success_http_code: '200'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Webhook payload that we expect from the user or VCS
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
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
// - in: path
//   name: hook
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
//     description: Successfully updated the webhook
//     schema:
//       type: string
//   '400':
//     description: The webhook was unable to be updated
//     schema:
//       type: string
//   '404':
//     description: The webhook was unable to be updated
//     schema:
//       type: string
//   '500':
//     description: The webhook was unable to be updated
//     schema:
//       type: string

// UpdateHook represents the API handler to update
// a webhook in the configured backend.
func UpdateHook(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	hook := c.Param("hook")

	logrus.Infof("Updating webhook %s/%s", r.GetFullName(), hook)

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for webhook %s/%s: %w", r.GetFullName(), hook, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHook(number, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get webhook %s/%d: %w", r.GetFullName(), number, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update webhook fields if provided
	if input.GetCreated() > 0 {
		// update created if set
		h.SetCreated(input.GetCreated())
	}

	if len(input.GetHost()) > 0 {
		// update host if set
		h.SetHost(input.GetHost())
	}

	if len(input.GetEvent()) > 0 {
		// update event if set
		h.SetEvent(input.GetEvent())
	}

	if len(input.GetBranch()) > 0 {
		// update branch if set
		h.SetBranch(input.GetBranch())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		h.SetError(input.GetError())
	}

	if len(input.GetStatus()) > 0 {
		// update status if set
		h.SetStatus(input.GetStatus())
	}

	if len(input.GetLink()) > 0 {
		// update link if set
		h.SetLink(input.GetLink())
	}

	// send API call to update the webhook
	err = database.FromContext(c).UpdateHook(h)
	if err != nil {
		retErr := fmt.Errorf("unable to update webhook %s/%d: %w", r.GetFullName(), h.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated user
	h, _ = database.FromContext(c).GetHook(h.GetNumber(), r)

	c.JSON(http.StatusOK, h)
}

// swagger:operation DELETE /api/v1/hooks/{org}/{repo}/{hook} deployment DeleteHooks
//
// Delete a webhook for the configured backend
//
// ---
// x-success_http_code: '200'
// x-incident_priority: P4
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
// - in: path
//   name: hook
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
//     description: Successfully updated the webhook
//     schema:
//       type: string
//   '400':
//     description: The webhook was unable to be updated
//     schema:
//       type: string
//   '404':
//     description: The webhook was unable to be updated
//     schema:
//       type: string
//   '500':
//     description: The webhook was unable to be updated
//     schema:
//       type: string

// DeleteHook represents the API handler to remove
// a webhook from the configured backend.
func DeleteHook(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	hook := c.Param("hook")

	logrus.Infof("Deleting webhook %s/%s", r.GetFullName(), hook)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHook(number, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get webhook %s/%d: %w", r.GetFullName(), number, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to remove the webhook
	err = database.FromContext(c).DeleteHook(h.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete webhook %s/%d: %w", r.GetFullName(), h.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Webhook %s/%d deleted", r.GetFullName(), h.GetNumber()))
}
