// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with service
package log

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps GetStepLog
//
// Get the logs for a step
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the logs for step
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"
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

// GetStepLog represents the API handler to get the logs for a step.
func GetStepLog(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	l.Debugf("reading logs for step %s", entry)

	// send API call to capture the step logs
	sl, err := database.FromContext(c).GetLogForStep(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, sl)
}
