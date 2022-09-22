// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package admin

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/admin/hooks admin AdminAllHooks
//
// Get all of the webhooks stored in the database
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved all hooks from the database
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Webhook"
//   '500':
//     description: Unable to retrieve all hooks
//     schema:
//       "$ref": "#/definitions/Error"

// AllHooks represents the API handler to
// captures all hooks stored in the database.
func AllHooks(c *gin.Context) {
	logrus.Info("Admin: reading all hooks")

	// send API call to capture all hooks
	r, err := database.FromContext(c).GetHookList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all hooks: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, r)
}

// swagger:operation PUT /api/v1/admin/hook admin AdminUpdateHook
//
// Update a hook in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing hook to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the hook in the database
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '404':
//     description: Unable to update the hook in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '501':
//     description: Unable to update the hook in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateHook represents the API handler to
// update any hook stored in the database.
func UpdateHook(c *gin.Context) {
	logrus.Info("Admin: updating hook in database")

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for hook %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the hook
	_, err = database.FromContext(c).UpdateHook(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update hook %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
