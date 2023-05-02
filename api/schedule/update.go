// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"fmt"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/router/middleware/repo"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/schedules/{org}/{repo}/{schedule} schedules UpdateSchedule
//
// Update a user for the configured backend
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
// - in: path
//   name: schedule
//   description: Name of the schedule
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the schedule to update
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
//     description: Unable to update the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to update the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the schedule
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateSchedule represents the API handler to update
// a schedule in the configured backend.
func UpdateSchedule(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	scheduleName := util.PathParameter(c, "schedule")

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"schedule": scheduleName,
		"repo":     r.GetName(),
		"org":      r.GetOrg(),
	}).Infof("updating schedule %s", scheduleName)

	// capture body from API request
	input := new(types.Schedule)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for schedule %s: %w", scheduleName, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the schedule
	s, err := database.FromContext(c).GetScheduleForRepo(r, scheduleName)
	if err != nil {
		retErr := fmt.Errorf("unable to get schedule %s: %w", scheduleName, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update schedule fields if provided
	if input.GetActive() {
		// update active if set to true
		s.SetActive(input.GetActive())
	}

	if input.GetName() != "" {
		// update name if defined
		s.SetName(input.GetName())
	}

	if input.GetEntry() != "" {
		// update entry if defined
		s.SetEntry(input.GetEntry())
	}

	// send API call to update the user
	err = database.FromContext(c).UpdateSchedule(s)
	if err != nil {
		retErr := fmt.Errorf("unable to update scheduled %s: %w", scheduleName, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated scheduled
	s, _ = database.FromContext(c).GetScheduleForRepo(r, scheduleName)

	c.JSON(http.StatusOK, s)
}
