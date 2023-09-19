// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with service
package log

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps UpdateStepLog
//
// Update the logs for a step
//
// ---
// deprecated: true
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
//   description: Payload containing the log to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Log"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the logs for step
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to update the logs for a step
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the logs for a step
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateStepLog represents the API handler to update
// the logs for a step in the configured backend.
func UpdateStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"step":  s.GetNumber(),
		"user":  u.GetName(),
	}).Infof("updating logs for step %s", entry)

	// send API call to capture the step logs
	l, err := database.FromContext(c).GetLogForStep(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// capture body from API request
	input := new(library.Log)

	err = c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update log fields if provided
	if len(input.GetData()) > 0 {
		// update data if set
		l.SetData(input.GetData())
	}

	// send API call to update the log
	err = database.FromContext(c).UpdateLog(ctx, l)
	if err != nil {
		retErr := fmt.Errorf("unable to update logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, nil)
}
