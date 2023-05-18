// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"net/http"

	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//
// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/services/{service} services GetService
//
// Get a service for a build in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: service
//   description: Name of the service
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the service
//     schema:
//       "$ref": "#/definitions/Service"
//   '400':
//     description: Unable to retrieve the service
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the service
//     schema:
//       "$ref": "#/definitions/Error"

// GetService represents the API handler to capture a
// service for a build from the configured backend.
func GetService(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"org":     o,
		"repo":    r.GetName(),
		"service": s.GetNumber(),
		"user":    u.GetName(),
	}).Infof("reading service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	c.JSON(http.StatusOK, s)
}
