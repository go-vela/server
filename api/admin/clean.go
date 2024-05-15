// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
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

	report := types.CleanReport{}

	logrus.Infof("platform admin %s: cleaning pending/running resources in database", u.GetName())

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

	if len(c.Query("before")) == 0 {
		retErr := fmt.Errorf("`before` query parameter is required")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture before query parameter, default to max build timeout
	before, err := strconv.ParseInt(c.Query("before"), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert before query parameter %s to int64: %w", c.Query("before"), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if cleanBuilds, _ := strconv.ParseBool(c.Query("builds")); cleanBuilds {
		// append pending approval to statuses if requested
		statuses := []string{constants.StatusRunning, constants.StatusPending}
		if cleanPendingApproval, _ := strconv.ParseBool(c.Query("pending_approval_builds")); cleanPendingApproval {
			statuses = append(statuses, constants.StatusPendingApproval)
		}

		// send API call to clean builds
		report.Builds, err = database.FromContext(c).CleanBuilds(ctx, msg, statuses, before)
		if err != nil {
			retErr := fmt.Errorf("unable to update builds: %w", err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		logrus.Infof("platform admin %s: cleaned %d builds in database", u.GetName(), report.Builds)

		// clean executables
		report.Executables, err = database.FromContext(c).CleanBuildExecutables(ctx)
		if err != nil {
			retErr := fmt.Errorf("%d builds cleaned. unable to clean build executables: %w", report.Builds, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		logrus.Infof("platform admin %s: cleaned %d build executables in database", u.GetName(), report.Executables)
	}

	if cleanServices, _ := strconv.ParseBool(c.Query("services")); cleanServices {
		// clean services
		report.Services, err = database.FromContext(c).CleanServices(ctx, msg, before)
		if err != nil {
			retErr := fmt.Errorf("%d builds cleaned. %d executables cleaned. unable to update services: %w", report.Builds, report.Executables, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		logrus.Infof("platform admin %s: cleaned %d services in database", u.GetName(), report.Services)
	}

	if cleanSteps, _ := strconv.ParseBool(c.Query("steps")); cleanSteps {
		// clean steps
		report.Steps, err = database.FromContext(c).CleanSteps(ctx, msg, before)
		if err != nil {
			retErr := fmt.Errorf("%d builds cleaned. %d executables cleaned. %d services cleaned. unable to update steps: %w", report.Builds, report.Executables, report.Services, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		logrus.Infof("platform admin %s: cleaned %d steps in database", u.GetName(), report.Steps)
	}

	c.JSON(http.StatusOK, report)
}
