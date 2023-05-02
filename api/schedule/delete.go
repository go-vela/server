// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"fmt"
	"github.com/go-vela/server/router/middleware/schedule"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/{schedule} schedules DeleteSchedule
//
// Delete a schedule in the configured backend
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the schedule
//     schema:
//       type: string
//   '500':
//     description: Unable to delete the schedule
//     schema:
//       "$ref": "#/definitions/Error"
//   '510':
//     description: Unable to delete the schedule
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteSchedule represents the API handler to remove
// a schedule from the configured backend.
func DeleteSchedule(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	s := schedule.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("deleting schedule %s", s.GetName())

	err := database.FromContext(c).DeleteSchedule(s)
	if err != nil {
		retErr := fmt.Errorf("unable to delete schedule %s: %w", s.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("schedule %s deleted", s.GetName()))
}
