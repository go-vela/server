// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedules

import (
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/schedule"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/schedules/{org}/{repo}/{schedule} schedules GetSchedule
//
// Get a schedule in the configured backend
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
//     description: Successfully retrieved the schedule
//     schema:
//       "$ref": "#/definitions/Schedule"

// GetSchedule represents the API handler to
// capture a schedule from the configured backend.
func GetSchedule(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	s := schedule.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"repo":     r.GetName(),
		"user":     u.GetName(),
		"schedule": s.GetName(),
	}).Infof("reading schedule %s", s.GetName())

	c.JSON(http.StatusOK, s)
}
