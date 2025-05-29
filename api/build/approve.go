// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/queue/models"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/approve builds ApproveBuild
//
// Approve a build to run
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
//   name: build
//   description: Build number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Request processed but build was skipped
//     schema:
//       type: string
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

// ApproveBuild represents the API handler to approve a build to run.
func ApproveBuild(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("approving build %d", b.GetID())

	// verify build is in correct status
	if !strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) {
		retErr := fmt.Errorf("unable to approve build %s/%d: build not in pending approval state", r.GetFullName(), b.GetNumber())
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// verify user is not the sender of the build
	if strings.EqualFold(u.GetName(), b.GetSender()) {
		retErr := fmt.Errorf("unable to approve build %s/%d: approver cannot be the sender of the build", r.GetFullName(), b.GetNumber())
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// set fields
	b.SetStatus(constants.StatusPending)
	b.SetApprovedAt(time.Now().Unix())
	b.SetApprovedBy(u.GetName())

	// update the build in the db
	_, err := database.FromContext(c).UpdateBuild(ctx, b)
	if err != nil {
		l.Errorf("failed to update build during publish to queue: %v", err)
	}

	l.Info("build updated - user approved build execution")

	// publish the build to the queue
	go Enqueue(
		context.WithoutCancel(ctx),
		queue.FromGinContext(c),
		database.FromContext(c),
		models.ToItem(b),
		b.GetRoute(),
	)

	c.JSON(http.StatusOK, fmt.Sprintf("Successfully approved build %s/%d", r.GetFullName(), b.GetNumber()))
}
