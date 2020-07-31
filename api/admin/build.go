// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package admin

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/admin/builds admin AdminAllBuilds
//
// Get all of the builds in the database
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved all builds from the database
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '500':
//     description: Unable to retrieve all builds from the database
//     schema:
//       type: string

// AllBuilds represents the API handler to
// captures all builds stored in the database.
func AllBuilds(c *gin.Context) {
	logrus.Info("Admin: reading all builds")

	// send API call to capture all builds
	b, err := database.FromContext(c).GetBuildList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all builds: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}

// swagger:operation PUT /api/v1/admin/build admin AdminUpdateBuild
//
// Update a build in the database
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully updated the build in the database
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '404':
//     description: Unable to update the build in the database
//     schema:
//       type: string
//   '500':
//     description: Unable to update the build in the database
//     schema:
//       type: string

// UpdateBuild represents the API handler to
// update any build stored in the database.
func UpdateBuild(c *gin.Context) {
	logrus.Info("Admin: updating build in database")

	// capture body from API request
	input := new(library.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for build %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the build
	err = database.FromContext(c).UpdateBuild(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update build %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
