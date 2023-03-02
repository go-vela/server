// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/initsteps steps CreateInitStep
//
// Create an initstep for a build
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
// - in: body
//   name: body
//   description: Payload containing the initstep to create
//   required: true
//   schema:
//     "$ref": "#/definitions/InitStep"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the initstep
//     type: json
//     schema:
//       "$ref": "#/definitions/InitStep"
//   '400':
//     description: Unable to create the initstep
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to create the initstep
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the initstep
//     schema:
//       "$ref": "#/definitions/Error"

// CreateInitStep represents the API handler to
// create an InitStep in the configured backend.
func CreateInitStep(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":   o,
		"repo":  r.GetName(),
		"build": b.GetNumber(),
		"user":  u.GetName(),
	})

	logger.Infof("creating new initstep for build %s", entry)

	// capture body from API request
	input := new(library.InitStep)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new initstep for build %s: %w", entry, err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in initstep object
	input.SetRepoID(r.GetID())
	input.SetBuildID(b.GetID())

	// send API call to create the initstep
	err = database.FromContext(c).CreateInitStep(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create initstep for build %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created initstep
	output, err := database.FromContext(c).GetInitStepForBuild(b, input.GetNumber())
	if err != nil {
		retErr := fmt.Errorf("unable to capture initstep %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, output)
}
