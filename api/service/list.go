// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/services services ListServices
//
// Get a list of all services for a build in the configured backend
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
//     description: Successfully retrieved the list of services
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Service"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of services
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of services
//     schema:
//       "$ref": "#/definitions/Error"

// ListServices represents the API handler to capture a list
// of services for a build from the configured backend.
func ListServices(c *gin.Context) {
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
	}).Infof("reading services for build %s", entry)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
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

	// send API call to capture the list of services for the build
	s, t, err := database.FromContext(c).ListServicesForBuild(b, map[string]interface{}{}, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get services for build %s: %w", entry, err)

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

	c.JSON(http.StatusOK, s)
}
