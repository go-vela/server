// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code
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

// swagger:operation GET /api/v1/admin/users admin AdminAllUsers
//
// Get all of the users in the database
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved all users from the database
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/User"
//   '500':
//     description: Unable to retrieve all users from the database
//     schema:
//       type: string

// AllUsers represents the API handler to
// captures all users stored in the database.
func AllUsers(c *gin.Context) {
	logrus.Info("Admin: reading all users")

	// send API call to capture all users
	u, err := database.FromContext(c).ListUsers()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all users: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, u)
}

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
	logrus.Info("Admin: updating user in database")

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
