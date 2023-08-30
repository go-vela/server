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

// swagger:operation POST /api/v1/users users CreateUser
//
// Create a user for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the user to create
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the user
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to create the user
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the user
//     schema:
//       "$ref": "#/definitions/Error"

// CreateUser represents the API handler to create
// a user in the configured backend.
func CreateUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new user: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("creating new user %s", input.GetName())

	// send API call to create the user
	user, err := database.FromContext(c).CreateUser(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create user: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, user)
}
