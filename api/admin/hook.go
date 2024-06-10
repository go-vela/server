// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation PUT /api/v1/admin/hook admin AdminUpdateHook
//
// Update a hook
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The hook object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the hook
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateHook represents the API handler to update a hook.
func UpdateHook(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()
	u := user.Retrieve(c)

	logger := logrus.WithFields(logrus.Fields{
		"ip":      util.EscapeValue(c.ClientIP()),
		"path":    util.EscapeValue(c.Request.URL.Path),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	})

	logrus.Debug("platform admin: updating hook")

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for hook %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"hook_id": input.GetID(),
	}).Debug("attempting to update hook")

	// send API call to update the hook
	h, err := database.FromContext(c).UpdateHook(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update hook %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"hook_id": h.GetID(),
	}).Info("hook updated")

	c.JSON(http.StatusOK, h)
}
