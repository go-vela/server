// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/hooks/{org}/{repo} webhook CreateHook
//
// Create a webhook for the configured backend
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
//   '201':
//     description: The webhook has been created
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: The webhook was unable to be created
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: The webhook was unable to be created
//     schema:
//       "$ref": "#/definitions/Error"

// CreateHook represents the API handler to create
// a webhook in the configured backend.
func CreateHook(c *gin.Context) {
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
	}).Infof("creating new hook for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new hook for repo %s: %w", r.GetFullName(), err)

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
	h, err := database.FromContext(c).CreateHook(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create hook for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, h)
}

// swagger:operation GET /api/v1/hooks/{org}/{repo} webhook GetHooks
//
// Retrieve the webhooks for the configured backend
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved webhooks
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Webhook"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve webhooks
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve webhooks
//     schema:
//       "$ref": "#/definitions/Error"

// GetHooks represents the API handler to capture a list
// of webhooks from the configured backend.
func GetHooks(c *gin.Context) {
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
	}).Infof("reading hooks for repo %s", r.GetFullName())

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

// swagger:operation GET /api/v1/hooks/{org}/{repo}/{hook} webhook GetHook
//
// Retrieve a webhook for the configured backend
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
// - in: path
//   name: hook
//   description: Number of the hook
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the webhook
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: Unable to retrieve the webhook
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the webhook
//     schema:
//       "$ref": "#/definitions/Error"

// GetHook represents the API handler to capture a
// webhook from the configured backend.
func GetHook(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	hook := util.PathParameter(c, "hook")

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), hook)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"hook": hook,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("reading hook %s", entry)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHook(number, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get hook %s: %w", entry, err)

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
// - in: path
//   name: hook
//   description: Number of the hook
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Webhook payload that we expect from the user or VCS
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the webhook
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: The webhook was unable to be updated
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: The webhook was unable to be updated
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: The webhook was unable to be updated
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateHook represents the API handler to update
// a webhook in the configured backend.
func UpdateHook(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	hook := util.PathParameter(c, "hook")

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), hook)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"hook": hook,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("updating hook %s", entry)

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for hook %s: %w", entry, err)

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
		retErr := fmt.Errorf("unable to get hook %s: %w", entry, err)

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
	h, err = database.FromContext(c).UpdateHook(h)
	if err != nil {
		retErr := fmt.Errorf("unable to update hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, h)
}

// swagger:operation DELETE /api/v1/hooks/{org}/{repo}/{hook} webhook DeleteHook
//
// Delete a webhook for the configured backend
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
// - in: path
//   name: hook
//   description: Number of the hook
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the webhook
//     schema:
//       type: string
//   '400':
//     description: The webhook was unable to be deleted
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: The webhook was unable to be deleted
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: The webhook was unable to be deleted
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteHook represents the API handler to remove
// a webhook from the configured backend.
func DeleteHook(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	hook := util.PathParameter(c, "hook")

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), hook)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"hook": hook,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("deleting hook %s", entry)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHook(number, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get hook %s: %w", hook, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to remove the webhook
	err = database.FromContext(c).DeleteHook(h.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete hook %s: %w", hook, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("hook %s deleted", entry))
}

// swagger:operation POST /api/v1/hooks/{org}/{repo}/{hook}/redeliver webhook RedeliverHook
//
// Redeliver a webhook from the SCM
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
// - in: path
//   name: hook
//   description: Number of the hook
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully redelivered the webhook
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: The webhook was unable to be redelivered
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: The webhook was unable to be redelivered
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: The webhook was unable to be redelivered
//     schema:
//       "$ref": "#/definitions/Error"

// RedeliverHook represents the API handler to redeliver
// a webhook from the SCM.

func RedeliverHook(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	hook := util.PathParameter(c, "hook")

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), hook)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"hook": hook,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("redelivering hook %s", entry)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHook(number, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get hook %s: %w", entry, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	err = scm.FromContext(c).RedeliverWebhook(c, u, r, h)
	if err != nil {
		retErr := fmt.Errorf("unable to redeliver hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("hook %s redelivered", entry))
}
