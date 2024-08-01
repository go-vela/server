// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/schedule"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/schedules/{org}/{repo}/{schedule} schedules UpdateSchedule
//
// Update a schedule
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
// - in: body
//   name: body
//   description: The schedule object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Schedule"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the schedule
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
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateSchedule represents the API handler to update a schedule.
func UpdateSchedule(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	s := schedule.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()
	scheduleName := util.PathParameter(c, "schedule")
	minimumFrequency := c.Value("scheduleminimumfrequency").(time.Duration)

	l.Debugf("updating schedule %s", scheduleName)

	// capture body from API request
	input := new(api.Schedule)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for schedule %s: %w", scheduleName, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update schedule fields if provided
	if input.Active != nil {
		// update active if set to true
		s.SetActive(input.GetActive())
	}

	if input.GetName() != "" {
		// update name if defined
		s.SetName(input.GetName())
	}

	if input.GetEntry() != "" {
		err = validateEntry(minimumFrequency, input.GetEntry())
		if err != nil {
			retErr := fmt.Errorf("schedule entry of %s is invalid: %w", input.GetEntry(), err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// update entry if defined
		s.SetEntry(input.GetEntry())
	}

	// set the updated by field using claims
	s.SetUpdatedBy(u.GetName())

	if input.GetBranch() != "" {
		s.SetBranch(input.GetBranch())
	}

	// update the schedule within the database
	s, err = database.FromContext(c).UpdateSchedule(ctx, s, true)
	if err != nil {
		retErr := fmt.Errorf("unable to update scheduled %s: %w", scheduleName, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
