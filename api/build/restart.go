// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build} builds RestartBuild
//
// Restart a build in the configured backend
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
//   description: Build number to restart
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
//     description: Successfully restarted the build
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to restart the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to restart the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to restart the build
//     schema:
//       "$ref": "#/definitions/Error"

// RestartBuild represents the API handler to restart an existing build in the configured backend.
func RestartBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	cl := claims.Retrieve(c)
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
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

	// a build that is in a pending approval state cannot be restarted
	if strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) {
		retErr := fmt.Errorf("unable to restart build %s/%d: cannot restart a build pending approval", r.GetFullName(), b.GetNumber())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// set sender to the user who initiated the restart and parent to the previous build
	b.SetSender(cl.Subject)
	b.SetParent(b.GetNumber())

	logger.Debugf("Generating queue items for build %s", entry)

	// restart form
	generatorForm := GeneratorForm{
		Build:    b,
		Repo:     r,
		Metadata: m,
		BaseErr:  "unable to restart build",
		Source:   "restart",
		Retries:  1,
	}

	// generate queue items
	_, _, items, err := GenerateQueueItems(
		c,
		generatorForm,
		database.FromContext(c),
		scm.FromContext(c),
		compiler.FromContext(c),
		queue.FromContext(c),
	)

	// error handling done in GenerateQueueItems
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, items.Build)

	// publish the build to the queue
	go PublishToQueue(
		ctx,
		queue.FromGinContext(c),
		database.FromContext(c),
		items.Build,
		items.Repo,
		items.User,
		items.Build.GetHost(),
	)
}
