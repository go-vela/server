// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adhocore/gronx"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/settings"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/schedules/{org}/{repo} schedules CreateSchedule
//
// Create a schedule
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
// - in: body
//   name: body
//   description: Schedule object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Schedule"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the schedule
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
//   '403':
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '409':
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"

// CreateSchedule represents the API handler to
// create a schedule.
func CreateSchedule(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()
	s := settings.FromContext(c)

	minimumFrequency := c.Value("scheduleminimumfrequency").(time.Duration)

	// capture body from API request
	input := new(api.Schedule)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new schedule: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure the entry is valid
	err = validateEntry(minimumFrequency, input.GetEntry())
	if err != nil {
		retErr := fmt.Errorf("schedule of %s with entry %s is invalid: %w", input.GetName(), input.GetEntry(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure schedule name is defined
	if input.GetName() == "" {
		util.HandleError(c, http.StatusBadRequest, fmt.Errorf("schedule name must be set"))

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Debugf("creating new schedule %s", input.GetName())

	// ensure repo is allowed to create new schedules
	if !util.CheckAllowlist(r, s.GetScheduleAllowlist()) {
		retErr := fmt.Errorf("unable to create schedule %s: %s is not on allowlist", input.GetName(), r.GetFullName())

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	schedule := new(api.Schedule)

	// update fields in schedule object
	schedule.SetCreatedBy(u.GetName())
	schedule.SetRepo(r)
	schedule.SetName(input.GetName())
	schedule.SetEntry(input.GetEntry())
	schedule.SetCreatedAt(time.Now().UTC().Unix())
	schedule.SetUpdatedAt(time.Now().UTC().Unix())
	schedule.SetUpdatedBy(u.GetName())

	if input.GetBranch() == "" {
		schedule.SetBranch(r.GetBranch())
	} else {
		schedule.SetBranch(input.GetBranch())
	}

	// set the active field based off the input provided
	if input.Active == nil {
		// default active field to true
		schedule.SetActive(true)
	} else {
		schedule.SetActive(input.GetActive())
	}

	// send API call to capture the schedule from the database
	dbSchedule, err := database.FromContext(c).GetScheduleForRepo(ctx, r, input.GetName())
	if err == nil && dbSchedule.GetActive() {
		retErr := fmt.Errorf("unable to create schedule: %s is already active", input.GetName())

		util.HandleError(c, http.StatusConflict, retErr)

		return
	}

	if !r.GetActive() {
		retErr := fmt.Errorf("unable to create schedule: %s repo %s is disabled", input.GetName(), r.GetFullName())

		util.HandleError(c, http.StatusConflict, retErr)

		return
	}

	// if the schedule exists but is inactive
	if dbSchedule.GetID() != 0 && !dbSchedule.GetActive() && input.GetActive() {
		// update the user who created the schedule
		dbSchedule.SetUpdatedBy(u.GetName())
		// activate the schedule
		dbSchedule.SetActive(true)

		// send API call to update the schedule
		schedule, err = database.FromContext(c).UpdateSchedule(ctx, dbSchedule, true)
		if err != nil {
			retErr := fmt.Errorf("unable to set schedule %s to active: %w", dbSchedule.GetName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	} else {
		// send API call to create the schedule
		schedule, err = database.FromContext(c).CreateSchedule(ctx, schedule)
		if err != nil {
			retErr := fmt.Errorf("unable to create new schedule %s: %w", r.GetName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	c.JSON(http.StatusCreated, schedule)
}

// validateEntry validates the entry for a minimum frequency.
func validateEntry(minimum time.Duration, entry string) error {
	gron := gronx.New()

	// check if expr is even valid
	valid := gron.IsValid(entry)
	if !valid {
		return fmt.Errorf("invalid entry of %s", entry)
	}

	// iterate 5 times through ticks in an effort to catch scalene entries
	tickForward := 5

	// start with now
	t := time.Now().UTC()

	for i := 0; i < tickForward; i++ {
		// check the previous occurrence of the entry
		prevTime, err := gronx.PrevTickBefore(entry, t, true)
		if err != nil {
			return err
		}

		// check the next occurrence of the entry
		nextTime, err := gronx.NextTickAfter(entry, t, false)
		if err != nil {
			return err
		}

		// ensure the time between previous and next schedule exceeds the minimum duration
		if nextTime.Sub(prevTime) < minimum {
			return fmt.Errorf("entry needs to occur less frequently than every %s", minimum)
		}

		t = nextTime
	}

	return nil
}
