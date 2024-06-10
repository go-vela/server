// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/users/{user} users UpdateUser
//
// Update a user
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: user
//   description: Name of the user
//   required: true
//   type: string
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
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateUser represents the API handler to update a user.
func UpdateUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	user := util.PathParameter(c, "user")
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Debugf("updating user %s", user)

	// capture body from API request
	input := new(types.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %s: %w", user, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the user
	u, err = database.FromContext(c).GetUserForName(ctx, user)
	if err != nil {
		retErr := fmt.Errorf("unable to get user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update user fields if provided
	if input.GetActive() {
		// update active if set to true
		u.SetActive(input.GetActive())
	}

	if input.GetAdmin() {
		// update admin if set to true
		u.SetAdmin(input.GetAdmin())
	}

	if input.Favorites != nil {
		// update favorites if set
		u.SetFavorites(input.GetFavorites())
	}

	if input.Dashboards != nil {
		// update dashboards if set
		u.SetDashboards(input.GetDashboards())
	}

	// send API call to update the user
	u, err = database.FromContext(c).UpdateUser(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", user, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, u)
}
