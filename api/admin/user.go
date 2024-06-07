// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
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
	logrus.Info("Admin: updating user in database")

	// capture middleware values
	ctx := c.Request.Context()

	// capture body from API request
	input := new(types.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to update the user
	u, err := database.FromContext(c).UpdateUser(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, u)
}
