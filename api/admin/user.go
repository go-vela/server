// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/admin/user admin AdminUpdateUser
//
// Update a user
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The user object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the user
//     schema:
//       "$ref": "#/definitions/User"
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

// UpdateUser represents the API handler to update a user.
func UpdateUser(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: updating user")

	// capture body from API request
	input := new(types.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"target_user_id": input.GetID(),
		"target_user":    util.EscapeValue(input.GetName()),
	}).Debug("platform admin: attempting to update user")

	// send API call to update the user
	tu, err := database.FromContext(c).UpdateUser(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"target_user_id": tu.GetID(),
		"target_user":    tu.GetName(),
	}).Info("platform admin: updated user")

	c.JSON(http.StatusOK, tu)
}
