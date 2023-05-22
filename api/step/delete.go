// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

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
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step} steps DeleteStep
//
// Delete a step for a build
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
//     description: Successfully deleted the step
//     schema:
//       type: string
//   '500':
//     description: Successfully deleted the step
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteStep represents the API handler to remove
// a step for a build from the configured backend.
func DeleteStep(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)

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
	}).Infof("deleting step %s", entry)

	// send API call to remove the step
	err := database.FromContext(c).DeleteStep(s)
	if err != nil {
		retErr := fmt.Errorf("unable to delete step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("step %s deleted", entry))
}
