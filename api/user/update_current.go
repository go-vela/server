// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/user users UpdateCurrentUser
//
// Update the current authenticated user in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the user to update
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the current user
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to update the current user
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to update the current user
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the current user
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateCurrentUser represents the API handler to capture and
// update the currently authenticated user from the configured backend.
func UpdateCurrentUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("updating current user %s", u.GetName())

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update user fields if provided
	if input.Favorites != nil {
		// update favorites if set
		u.SetFavorites(input.GetFavorites())
	}

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated user
	u, err = database.FromContext(c).GetUserForName(u.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to get updated user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	c.JSON(http.StatusOK, u)
}
