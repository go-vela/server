// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/initsteps steps ListInitSteps
//
// List initsteps for a build
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
// - in: query
//   name: page
//   description: The page of results to retrieve
//   type: integer
//   default: 1
// - in: query
//   name: per_page
//   description: How many results per page to return
//   type: integer
//   maximum: 100
//   default: 10
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the initsteps
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/InitStep"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of initsteps
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of initsteps
//     schema:
//       "$ref": "#/definitions/Error"

// ListInitSteps represents the API handler to capture a list
// of initsteps for a repo from the configured backend.
func ListInitSteps(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	b := build.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":   o,
		"repo":  r.GetName(),
		"build": b.GetNumber(),
		"user":  u.GetName(),
	}).Infof("listing initsteps for repo %s", entry)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		//nolint:lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to convert page query parameter for build %s: %w", entry, err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for build %s: %w", entry, err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of InitStep logs for the build
	output, t, err := database.FromContext(c).ListInitStepsForBuild(b, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get initsteps for build %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := api.Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, output)
}
