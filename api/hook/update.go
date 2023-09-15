// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

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
	ctx := c.Request.Context()

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
	h, err := database.FromContext(c).GetHookForRepo(ctx, r, number)
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
	h, err = database.FromContext(c).UpdateHook(ctx, h)
	if err != nil {
		retErr := fmt.Errorf("unable to update hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, h)
}
