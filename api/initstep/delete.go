// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"fmt"
	"github.com/go-vela/server/router/middleware/initstep"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/initsteps/{step} steps DeleteInitStep
//
// Delete an initstep for a build
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
//   name: initstep
//   description: InitStep number
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the initstep
//     schema:
//       type: string
//   '500':
//     description: Unable to delete the initstep
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteInitStep represents the API handler to remove
// an InitStep for a repo from the configured backend.
func DeleteInitStep(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	i := initstep.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), i.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"repo":     r.GetName(),
		"build":    b.GetNumber(),
		"initstep": i.GetNumber(),
		"user":     u.GetName(),
	}).Infof("deleting initstep %s", entry)

	// send API call to remove the InitStep
	err := database.FromContext(c).DeleteInitStep(i)
	if err != nil {
		retErr := fmt.Errorf("unable to delete initstep %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("initstep %s deleted", entry))
}
