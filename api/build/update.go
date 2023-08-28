// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build} builds UpdateBuild
//
// Updates a build in the configured backend
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
//   description: Build number to update
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing the build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the build
//     schema:
//       "$ref": "#/definitions/Build"
//   '404':
//     description: Unable to update the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the build
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateBuild represents the API handler to update
// a build for a repo in the configured backend.
func UpdateBuild(c *gin.Context) {
	// capture middleware values
	cl := claims.Retrieve(c)
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  cl.Subject,
	}).Infof("updating build %s", entry)

	// capture body from API request
	input := new(library.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for build %s: %w", entry, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update build fields if provided
	if len(input.GetStatus()) > 0 {
		// update status if set
		b.SetStatus(input.GetStatus())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		b.SetError(input.GetError())
	}

	if input.GetEnqueued() > 0 {
		// update enqueued if set
		b.SetEnqueued(input.GetEnqueued())
	}

	if input.GetStarted() > 0 {
		// update started if set
		b.SetStarted(input.GetStarted())
	}

	if input.GetFinished() > 0 {
		// update finished if set
		b.SetFinished(input.GetFinished())
	}

	if len(input.GetTitle()) > 0 {
		// update title if set
		b.SetTitle(input.GetTitle())
	}

	if len(input.GetMessage()) > 0 {
		// update message if set
		b.SetMessage(input.GetMessage())
	}

	if len(input.GetHost()) > 0 {
		// update host if set
		b.SetHost(input.GetHost())
	}

	if len(input.GetRuntime()) > 0 {
		// update runtime if set
		b.SetRuntime(input.GetRuntime())
	}

	if len(input.GetDistribution()) > 0 {
		// update distribution if set
		b.SetDistribution(input.GetDistribution())
	}

	// send API call to update the build
	b, err = database.FromContext(c).UpdateBuild(ctx, b)
	if err != nil {
		retErr := fmt.Errorf("unable to update build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)

	// check if the build is in a "final" state
	if b.GetStatus() == constants.StatusSuccess ||
		b.GetStatus() == constants.StatusFailure ||
		b.GetStatus() == constants.StatusCanceled ||
		b.GetStatus() == constants.StatusKilled ||
		b.GetStatus() == constants.StatusError {
		// send API call to capture the repo owner
		u, err := database.FromContext(c).GetUser(ctx, r.GetUserID())
		if err != nil {
			logrus.Errorf("unable to get owner for build %s: %v", entry, err)
		}

		// send API call to set the status on the commit
		err = scm.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for build %s: %v", entry, err)
		}
	}
}
