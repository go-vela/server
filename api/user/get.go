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
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/users/{user} users GetUser
//
// Retrieve a user for the configured backend
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
//     description: Successfully retrieved the user
//     schema:
//       "$ref": "#/definitions/User"
//   '404':
//     description: Unable to retrieve the user
//     schema:
//       "$ref": "#/definitions/Error"

// GetUser represents the API handler to capture a
// user from the configured backend.
func GetUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	user := util.PathParameter(c, "user")

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("reading user %s", user)

	// send API call to capture the user
	u, err := database.FromContext(c).GetUserForName(user)
	if err != nil {
		retErr := fmt.Errorf("unable to get user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	c.JSON(http.StatusOK, u)
}
