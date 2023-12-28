// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/admin/clean admin AdminCleanResources
//
// Update pending build resources to error status before a given time
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: before
//   description: filter pending resources created before a certain time
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
//     type: string
//   '400':
//     description: Unable to update resources — bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unable to update resources — unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update resources
//     schema:
//       "$ref": "#/definitions/Error"

// CleanResources represents the API handler to
// update any user stored in the database.
func CleanResources(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	logrus.Infof("platform admin %s: updating pending resources in database", u.GetName())

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
		msg = *input.Message
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

	logrus.Infof("platform admin %s: cleaned %d builds in database", u.GetName(), builds)

	// clean executables
	executables, err := database.FromContext(c).CleanBuildExecutables(ctx)
	if err != nil {
		retErr := fmt.Errorf("%d builds cleaned. unable to clean build executables: %w", builds, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// clean services
	services, err := database.FromContext(c).CleanServices(ctx, msg, before)
	if err != nil {
		retErr := fmt.Errorf("%d builds cleaned. %d executables cleaned. unable to update services: %w", builds, executables, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logrus.Infof("platform admin %s: cleaned %d services in database", u.GetName(), services)

	// clean steps
	steps, err := database.FromContext(c).CleanSteps(msg, before)
	if err != nil {
		retErr := fmt.Errorf("%d builds cleaned. %d executables cleaned. %d services cleaned. unable to update steps: %w", builds, executables, services, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logrus.Infof("platform admin %s: cleaned %d steps in database", u.GetName(), steps)

	c.JSON(http.StatusOK, fmt.Sprintf("%d builds cleaned. %d executables cleaned. %d services cleaned. %d steps cleaned.", builds, executables, services, steps))
}
