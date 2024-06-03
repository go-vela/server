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
//   '401':
//     description: Unauthorized to update the hook in the database
//     schema:
//       "$ref": "#/definitions/Error
//   '400':
//     description: Unable to update the hook in the database - bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the hook in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateHook represents the API handler to
// update any hook stored in the database.
func UpdateHook(c *gin.Context) {
	logrus.Debug("platform admin: updating hook")

	// capture middleware values
	ctx := c.Request.Context()
	u := user.Retrieve(c)

	logger := logrus.WithFields(logrus.Fields{
		"ip":      util.EscapeValue(c.ClientIP()),
		"path":    util.EscapeValue(c.Request.URL.Path),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	})

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
	}).Info("attempting to update hook")

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
