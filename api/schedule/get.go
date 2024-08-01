// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/schedule"
)

// swagger:operation GET /api/v1/schedules/{org}/{repo}/{schedule} schedules GetSchedule
//
// Get a schedule
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
//   name: schedule
//   description: Name of the schedule
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the schedule
//     schema:
//       "$ref": "#/definitions/Schedule"
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

// GetSchedule represents the API handler to get a schedule.
func GetSchedule(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	s := schedule.Retrieve(c)

	l.Debugf("reading schedule %s", s.GetName())

	c.JSON(http.StatusOK, s)
}
