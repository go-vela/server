// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

// swagger:operation POST /api/v1/schedules/{org}/{repo} schedules CreateSchedule
//
// Create a schedule in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the schedule to create
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
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '403':
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '409':
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to create the schedule
//     schema:
//       "$ref": "#/definitions/Error"

// CreateSchedule represents the API handler to
// create a schedule in the configured backend.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func CreateSchedule(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	r := repo.Retrieve(c)
	allowlist := c.Value("allowlistschedule").([]string)

	// capture body from API request
	input := new(types.Schedule)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new schedule: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// TODO: add code to validate the input.Entry matches what we allow

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("creating new schedule %s", input.GetName())

	// ensure repo is allowed to be activated
	if !util.CheckAllowlist(r, allowlist) {
		retErr := fmt.Errorf("unable to create schedule %s: %s is not on allowlist", input.GetName(), r.GetFullName())

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	s := new(types.Schedule)

	// update fields in repo object
	r.SetUserID(u.GetID())

	// set the active field based off the input provided
	if input.Active == nil {
		// default active field to true
		s.SetActive(true)
	} else {
		s.SetActive(input.GetActive())
	}

	// send API call to capture the schedule from the database
	dbSchedule, err := database.FromContext(c).GetScheduleForRepo(r, input.GetName())
	if err == nil && dbSchedule.GetActive() {
		retErr := fmt.Errorf("unable to create schedule: %s is already active", input.GetName())

		util.HandleError(c, http.StatusConflict, retErr)

		return
	}

	// send API call to capture the repo from the database
	dbRepo, err := database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())
	if err == nil && !dbRepo.GetActive() {
		retErr := fmt.Errorf("unable to create schedule: %s repo %s is disabled", input.GetName(), dbRepo.GetFullName())

		util.HandleError(c, http.StatusConflict, retErr)

		return
	}

	// if the repo exists but is inactive
	if len(dbRepo.GetOrg()) > 0 && !dbSchedule.GetActive() && input.GetActive() {
		// update the repo owner
		dbSchedule.SetCreatedBy(u.GetName())
		// activate the schedule
		dbSchedule.SetActive(true)

		// send API call to update the repo
		err = database.FromContext(c).UpdateSchedule(dbSchedule)
		if err != nil {
			retErr := fmt.Errorf("unable to set schedule %s to active: %w", dbSchedule.GetName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture the updated repo
		s, _ = database.FromContext(c).GetScheduleForRepo(r, dbSchedule.GetName())
	} else {
		// send API call to create the repo
		err = database.FromContext(c).CreateSchedule(s)
		if err != nil {
			retErr := fmt.Errorf("unable to create new schedule %s: %w", r.GetName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture the created repo
		s, _ = database.FromContext(c).GetScheduleForRepo(r, dbSchedule.GetName())
	}

	c.JSON(http.StatusCreated, s)
}
