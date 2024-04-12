// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds builds CreateBuild
//
// Create a build in the configured backend
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
// - in: body
//   name: body
//   description: Payload containing the build to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully received the webhook but build was skipped
//     schema:
//       type: string
//   '201':
//     description: Successfully created the build from webhook
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Malformed webhook payload or improper pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Repository owner does not have proper access
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to receive the webhook
//     schema:
//       "$ref": "#/definitions/Error"
//   '429':
//     description: Concurrent build limit reached for repository
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to receive the webhook or internal error while processing
//     schema:
//       "$ref": "#/definitions/Error"

// CreateBuild represents the API handler to create a build in the configured backend.
func CreateBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
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

	logger.Infof("creating new build for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new build for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// verify the build has a valid event and the repo allows that event type
	if !r.GetAllowEvents().Allowed(input.GetEvent(), input.GetEventAction()) {
		retErr := fmt.Errorf("unable to create new build: %s does not have %s events enabled", r.GetFullName(), input.GetEvent())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create config
	config := CompileAndPublishConfig{
		Build:    input,
		Repo:     r,
		Metadata: m,
		BaseErr:  "unable to create build",
		Source:   "create",
		Retries:  1,
	}

	_, item, code, err := CompileAndPublish(
		c,
		config,
		database.FromContext(c),
		scm.FromContext(c),
		compiler.FromContext(c),
		queue.FromContext(c),
	)

	// check if build was skipped
	if err != nil && code == http.StatusOK {
		c.JSON(http.StatusOK, err.Error())

		return
	}

	if err != nil {
		util.HandleError(c, code, err)
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
