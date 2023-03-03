// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/initstep"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/initsteps/{step} initsteps GetInitStep
//
// Retrieve an initstep for a build
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
//   name: initstep
//   description: InitStep number
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the initstep
//     type: json
//     schema:
//       "$ref": "#/definitions/InitStep"

// GetInitStep represents the API handler to capture
// an InitStep for a repo from the configured backend.
func GetInitStep(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	i := initstep.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"repo":     r.GetName(),
		"build":    b.GetNumber(),
		"initstep": i.GetNumber(),
		"user":     u.GetName(),
	}).Infof("reading initstep %s/%d/%d", r.GetFullName(), b.GetNumber(), i.GetNumber())

	c.JSON(http.StatusOK, i)
}
