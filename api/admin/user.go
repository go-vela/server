// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/admin/user admin AdminUpdateUser
//
// Update a user in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing user to update
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the user in the database
//     schema:
//       "$ref": "#/definitions/User"
//   '404':
//     description: Unable to update the user in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the user in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateUser represents the API handler to
// update any user stored in the database.
func UpdateUser(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()

	// capture body from API request
	input := new(types.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

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
