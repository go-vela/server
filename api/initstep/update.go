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
	"github.com/go-vela/server/router/middleware/initstep"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/initsteps/{step} initsteps UpdateInitStep
//
// Update an initstep for a build
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
// - in: body
//   name: body
//   description: Payload containing the initstep to update
//   required: true
//   schema:
//     "$ref": "#/definitions/InitStep"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the initstep
//     schema:
//       "$ref": "#/definitions/InitStep"
//   '404':
//     description: Unable to update the initstep
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the initstep
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateInitStep represents the API handler to update
// an InitStep for a repo in the configured backend.
func UpdateInitStep(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	i := initstep.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), i.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"repo":     r.GetName(),
		"build":    b.GetNumber(),
		"initstep": i.GetNumber(),
		"user":     u.GetName(),
	}).Infof("updating initstep %s", entry)

	// capture body from API request
	input := new(library.InitStep)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for initstep %s: %w", entry, err)
		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// check if the Reporter field in the initstep was provided
	if len(input.GetReporter()) > 0 {
		// update the Reporter field
		i.SetReporter(input.GetReporter())
	}

	// check if the Name field in the initstep was provided
	if len(input.GetName()) > 0 {
		// update the Name field
		i.SetName(input.GetName())
	}

	// check if the Mimetype field in the initstep was provided
	if len(input.GetMimetype()) > 0 {
		// update the Mimetype field
		i.SetMimetype(input.GetMimetype())
	}

	// send API call to update the initstep
	err = database.FromContext(c).UpdateInitStep(i)
	if err != nil {
		retErr := fmt.Errorf("unable to update initstep %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated initstep
	i, err = database.FromContext(c).GetInitStepForBuild(b, i.GetNumber())
	if err != nil {
		retErr := fmt.Errorf("unable to capture initstep %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, i)
}
