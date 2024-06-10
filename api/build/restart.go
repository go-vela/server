// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
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
	cl := claims.Retrieve(c)
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	scm := scm.FromContext(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	})

	logger.Debugf("restarting build %d", b.GetNumber())

	// a build that is in a pending approval state cannot be restarted
	if strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) {
		retErr := fmt.Errorf("unable to restart build %s/%d: cannot restart a build pending approval", r.GetFullName(), b.GetNumber())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
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

	logger.Debugf("Generating queue items for build %s", entry)

	// restart form
	config := CompileAndPublishConfig{
		Build:    b,
		Metadata: m,
		BaseErr:  "unable to restart build",
		Source:   "restart",
		Retries:  1,
	}

	// generate queue items
	_, item, code, err := CompileAndPublish(
		c,
		config,
		database.FromContext(c),
		scm,
		compiler.FromContext(c),
		queue.FromContext(c),
	)
	if err != nil {
		util.HandleError(c, code, err)

		return
	}

	c.JSON(http.StatusCreated, item.Build)

	// publish the build to the queue
	go Enqueue(
		ctx,
		queue.FromGinContext(c),
		database.FromContext(c),
		item,
		item.Build.GetHost(),
	)
}
