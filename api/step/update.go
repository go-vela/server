// SPDX-License-Identifier: Apache-2.0

package step

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step} steps UpdateStep
//
// Update a step for a build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: step
//   description: Step number
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: The step object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Step"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the step
//     schema:
//       "$ref": "#/definitions/Step"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateStep represents the API handler to update
// a step for a build.
func UpdateStep(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	l.Debugf("updating step %s", entry)

	// capture body from API request
	input := new(types.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update step fields if provided
	if len(input.GetStatus()) > 0 {
		// update status if set
		s.SetStatus(input.GetStatus())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		s.SetError(input.GetError())
	}

	if input.GetExitCode() > 0 {
		// update exit_code if set
		s.SetExitCode(input.GetExitCode())
	}

	if input.GetStarted() > 0 {
		// update started if set
		s.SetStarted(input.GetStarted())
	}

	if input.GetFinished() > 0 {
		// update finished if set
		s.SetFinished(input.GetFinished())
	}

	if len(input.GetHost()) > 0 {
		// update host if set
		s.SetHost(input.GetHost())
	}

	if len(input.GetRuntime()) > 0 {
		// update runtime if set
		s.SetRuntime(input.GetRuntime())
	}

	if len(input.GetDistribution()) > 0 {
		// update distribution if set
		s.SetDistribution(input.GetDistribution())
	}

	// send API call to update the step
	s, err = database.FromContext(c).UpdateStep(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to update step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)

	// check if the build is in a "final" state
	// and if build is not a scheduled event
	if (s.GetStatus() == constants.StatusSuccess ||
		s.GetStatus() == constants.StatusFailure ||
		s.GetStatus() == constants.StatusCanceled ||
		s.GetStatus() == constants.StatusKilled ||
		s.GetStatus() == constants.StatusError) &&
		(b.GetEvent() != constants.EventSchedule) &&
		(len(s.GetReportAs()) > 0) {
		// send API call to set the status on the commit
		err = scm.FromContext(c).StepStatus(ctx, r.GetOwner(), b, s, r.GetOrg(), r.GetName())
		if err != nil {
			l.Errorf("unable to set commit status for build %s: %v", entry, err)
		}
	}
}
