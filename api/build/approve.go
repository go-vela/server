// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/approve builds ApproveBuild
//
// Sign off on a build to run from an outside contributor
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
//   description: Build number to retrieve
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Request processed but build was skipped
//     schema:
//       type: string
//   '201':
//     description: Successfully created the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to create the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to create the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the build
//     schema:
//       "$ref": "#/definitions/Error"

// CreateBuild represents the API handler to approve a build to run in the configured backend.
func ApproveBuild(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	})

	if !strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) {
		retErr := fmt.Errorf("unable to approve build %s/%d: build not in pending approval state", r.GetFullName(), b.GetNumber())
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	logger.Debugf("user %s approved build %s/%d for execution", u.GetName(), r.GetFullName(), b.GetNumber())

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(ctx, r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	b.SetStatus(constants.StatusPending)

	// update the build in the db
	_, err = database.FromContext(c).UpdateBuild(ctx, b)
	if err != nil {
		logrus.Errorf("Failed to update build %d during publish to queue for %s: %v", b.GetNumber(), r.GetFullName(), err)
	}

	// publish the build to the queue
	go PublishToQueue(
		ctx,
		queue.FromGinContext(c),
		database.FromContext(c),
		b,
		r,
		u,
		b.GetHost(),
	)

	c.JSON(http.StatusOK, fmt.Sprintf("Successfully approved build %s/%d", r.GetFullName(), b.GetNumber()))
}
