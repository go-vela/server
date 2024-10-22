// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/admin/clean admin AdminCleanResources
//
// Update stale build resources to error status
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: before
//   description: Filter stale resources created before a certain time
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing error message
//   required: true
//   schema:
//     "$ref": "#/definitions/Error"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated pending resources with error message
//     schema:
//       type: string
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// CleanResources represents the API handler to update stale resources.
func CleanResources(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: cleaning resources")

	// default error message
	msg := "build cleaned by platform admin"

	// capture body from API request
	input := new(types.Error)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for error message: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// if a message is provided, set the error message to that
	if input.Message != nil {
		msg = util.EscapeValue(*input.Message)
	}

	// capture before query parameter, default to max build timeout
	before, err := strconv.ParseInt(c.DefaultQuery("before", fmt.Sprint((time.Now().Add(-(time.Minute * (constants.BuildTimeoutMax + 5)))).Unix())), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert before query parameter %s to int64: %w", c.Query("before"), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to clean builds
	builds, err := database.FromContext(c).CleanBuilds(ctx, msg, before)
	if err != nil {
		retErr := fmt.Errorf("unable to update builds: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Debugf("platform admin: cleaned %d builds", builds)

	// clean executables
	executables, err := database.FromContext(c).CleanBuildExecutables(ctx)
	if err != nil {
		retErr := fmt.Errorf("%d builds cleaned. unable to clean build executables: %w", builds, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Debugf("platform admin: cleaned %d executables", executables)

	// clean services
	services, err := database.FromContext(c).CleanServices(ctx, msg, before)
	if err != nil {
		retErr := fmt.Errorf("%d builds cleaned. %d executables cleaned. unable to update services: %w", builds, executables, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Debugf("platform admin: cleaned %d services", services)

	// clean steps
	steps, err := database.FromContext(c).CleanSteps(ctx, msg, before)
	if err != nil {
		retErr := fmt.Errorf("%d builds cleaned. %d executables cleaned. %d services cleaned. unable to update steps: %w", builds, executables, services, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Debugf("platform admin: cleaned %d steps", steps)

	c.JSON(http.StatusOK, fmt.Sprintf("%d builds cleaned. %d executables cleaned. %d services cleaned. %d steps cleaned.", builds, executables, services, steps))
}
