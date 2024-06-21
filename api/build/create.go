// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds builds CreateBuild
//
// Create a build
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
// - in: body
//   name: body
//   description: Build object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
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

// CreateBuild represents the API handler to create a build.
func CreateBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("creating new build for repo %s", r.GetFullName())

	// capture body from API request
	input := new(types.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new build for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	input.SetRepo(r)

	// verify the build has a valid event and the repo allows that event type
	if !r.GetAllowEvents().Allowed(input.GetEvent(), input.GetEventAction()) {
		retErr := fmt.Errorf("unable to create new build: %s does not have %s events enabled", r.GetFullName(), input.GetEvent())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create config
	config := CompileAndPublishConfig{
		Build:    input,
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

		return
	}

	l.WithFields(logrus.Fields{
		"build":    item.Build.GetNumber(),
		"build_id": item.Build.GetID(),
	}).Info("build created")

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
