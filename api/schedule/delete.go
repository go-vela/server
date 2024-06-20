// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/schedule"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/{schedule} schedules DeleteSchedule
//
// Delete a schedule
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
//     description: Successfully deleted the schedule
//     schema:
//       type: string
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

// DeleteSchedule represents the API handler to remove a schedule.
func DeleteSchedule(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	s := schedule.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("deleting schedule %s", s.GetName())

	err := database.FromContext(c).DeleteSchedule(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to delete schedule %s: %w", s.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("schedule %s deleted", s.GetName()))
}
