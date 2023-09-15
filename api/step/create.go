// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/steps steps CreateStep
//
// Create a step for a build
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
//   description: Payload containing the step to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Step"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the step
//     schema:
//       "$ref": "#/definitions/Step"
//   '400':
//     description: Unable to create the step
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the step
//     schema:
//       "$ref": "#/definitions/Error"

// CreateStep represents the API handler to create
// a step for a build in the configured backend.
func CreateStep(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("creating new step for build %s", entry)

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new step for build %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in step object
	input.SetRepoID(r.GetID())
	input.SetBuildID(b.GetID())

	if len(input.GetStatus()) == 0 {
		input.SetStatus(constants.StatusPending)
	}

	if input.GetCreated() == 0 {
		input.SetCreated(time.Now().UTC().Unix())
	}

	// send API call to create the step
	s, err := database.FromContext(c).CreateStep(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create step for build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, s)
}
