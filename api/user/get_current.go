// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/user users GetCurrentUser
//
// Retrieve the current authenticated user from the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the current user
//     schema:
//       "$ref": "#/definitions/User"

// GetCurrentUser represents the API handler to capture the
// currently authenticated user from the configured backend.
func GetCurrentUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("reading current user %s", u.GetName())

	c.JSON(http.StatusOK, u)
}
