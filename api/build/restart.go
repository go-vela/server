// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/settings"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build} builds RestartBuild
//
// Restart a build
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
//     description: Successfully received request but build was skipped
//     schema:
//       type: string
//   '201':
//     description: Successfully created the build from request
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
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
//   '429':
//     description: Concurrent build limit reached for repository
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// RestartBuild represents the API handler to restart an existing build.
func RestartBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	cl := claims.Retrieve(c)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	scm := scm.FromContext(c)
	ps := settings.FromContext(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	l.Debugf("restarting build %d", b.GetNumber())

	// a build that is in a pending approval state cannot be restarted
	if strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) {
		retErr := fmt.Errorf("unable to restart build %s/%d: cannot restart a build pending approval", r.GetFullName(), b.GetNumber())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// deny restarts for merge queue builds. Users should enqueue a new build instead.
	if b.GetEvent() == constants.EventMergeGroup {
		retErr := fmt.Errorf("unable to restart build %s/%d: cannot restart a build in a merge group. Enqueue a new build instead", r.GetFullName(), b.GetNumber())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// check to see if queue has reached configured capacity to allow restarts for pending builds
	if b.GetStatus() == constants.StatusPending && ps.GetQueueRestartLimit() > 0 {
		// check length of specified route for the build
		queueLength, err := queue.FromContext(c).RouteLength(ctx, b.GetRoute())
		if err != nil {
			util.HandleError(c, http.StatusInternalServerError, fmt.Errorf("unable to get queue length for %s: %w", b.GetRoute(), err))

			return
		}

		if queueLength >= int64(ps.GetQueueRestartLimit()) {
			retErr := fmt.Errorf("unable to restart build %s: queue length %d exceeds configured limit %d, please wait for the queue to decrease in size before retrying", entry, queueLength, ps.GetQueueRestartLimit())

			util.HandleError(c, http.StatusTooManyRequests, retErr)

			return
		}
	}

	// set sender to the user who initiated the restart
	b.SetSender(cl.Subject)

	// fetch scm user id
	senderID, err := scm.GetUserID(ctx, u.GetName(), r.GetOwner().GetToken())
	if err != nil {
		retErr := fmt.Errorf("unable to get SCM user id for %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	b.SetSenderSCMID(senderID)

	// parent to the previous build
	b.SetParent(b.GetNumber())

	l.Debugf("generating queue items for build %s", entry)

	// restart form
	config := CompileAndPublishConfig{
		Build:    b,
		Metadata: m,
		BaseErr:  "unable to restart build",
		Source:   "restart",
		Retries:  1,
	}

	// generate queue items
	p, item, code, err := CompileAndPublish(
		c,
		config,
		database.FromContext(c),
		cache.FromContext(c),
		scm,
		compiler.FromContext(c),
		queue.FromContext(c),
	)
	if err != nil {
		util.HandleError(c, code, err)

		return
	}

	// determine whether or not to send compiled build to queue
	shouldEnqueue, err := ShouldEnqueue(c, l, item.Build, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	if shouldEnqueue {
		// send API call to set the status on the commit
		err := scm.Status(c.Request.Context(), b, r.GetOrg(), r.GetName(), p.Token)
		if err != nil {
			l.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
		}

		// publish the build to the queue
		go Enqueue(
			context.WithoutCancel(c.Request.Context()),
			queue.FromGinContext(c),
			database.FromContext(c),
			item,
			item.Build.GetRoute(),
		)
	} else {
		err := GatekeepBuild(c, item.Build, item.Build.GetRepo(), p.Token)
		if err != nil {
			util.HandleError(c, http.StatusInternalServerError, err)

			return
		}
	}

	l.WithFields(logrus.Fields{
		"new_build":    item.Build.GetNumber(),
		"new_build_id": item.Build.GetID(),
	}).Info("build created via restart")

	c.JSON(http.StatusCreated, item.Build)
}
