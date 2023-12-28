// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/users/{user} users DeleteUser
//
// Delete a user for the configured backend
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted of user
//     schema:
//       type: string
//   '404':
//     description: Unable to delete user
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to delete user
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteUser represents the API handler to remove
// a user from the configured backend.
func DeleteUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	user := util.PathParameter(c, "user")
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("deleting user %s", user)

	// send API call to capture the user
	u, err := database.FromContext(c).GetUserForName(ctx, user)
	if err != nil {
		retErr := fmt.Errorf("unable to get user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to remove the user
	err = database.FromContext(c).DeleteUser(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to delete user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("user %s deleted", u.GetName()))
}
