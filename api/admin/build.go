// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/admin/builds/:id admin AdminSingleBuild
//
// Get a single build in the database by ID
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved build from the database
//     schema:
//       "$ref": "#/definitions/Build"
//   '500':
//     description: Unable to retrieve build from the database
//     schema:
//       "$ref": "#/definitions/Error"

// SingleBuild represents the API handler to
// capture a single build stored in the database by id.
func SingleBuild(c *gin.Context) {
	// Logger
	logrus.Info("Admin: reading build")

	// Parse build ID from path
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		retErr := fmt.Errorf("unable to parse build id: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the build
	b, err := database.FromContext(c).GetBuildByID(id)
	if err != nil {
		retErr := fmt.Errorf("unable to capture build: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}

// swagger:operation GET /api/v1/admin/builds/queue admin AllBuildsQueue
//
// Get all of the running and pending builds in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: after
//   description: Unix timestamp to limit builds returned
//   required: false
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved all running and pending builds from the database
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/BuildQueue"
//   '500':
//     description: Unable to retrieve all running and pending builds from the database
//     schema:
//       "$ref": "#/definitions/Error"

// AllBuildsQueue represents the API handler to
// captures all running and pending builds stored in the database.
func AllBuildsQueue(c *gin.Context) {
	logrus.Info("Admin: reading running and pending builds")

	// default timestamp to 24 hours ago if user did not provide it as query parameter
	after := c.DefaultQuery("after", strconv.FormatInt(time.Now().UTC().Add(-24*time.Hour).Unix(), 10))

	// send API call to capture pending and running builds
	b, err := database.FromContext(c).GetPendingAndRunningBuilds(after)
	if err != nil {
		retErr := fmt.Errorf("unable to capture all running and pending builds: %w", err)

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
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the build in the database
//     schema:
//       "$ref": "#/definitions/Build"
//   '404':
//     description: Unable to update the build in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the build in the database
//     schema:
//       "$ref": "#/definitions/Error"

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
