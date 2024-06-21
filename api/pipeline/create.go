// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation POST /api/v1/pipelines/{org}/{repo} pipelines CreatePipeline
//
// Create a pipeline
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
//   description: Pipeline object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Pipeline"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/Pipeline"
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

// CreatePipeline represents the API handler to
// create a pipeline.
func CreatePipeline(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("creating new pipeline for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Pipeline)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new build for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in pipeline object
	input.SetRepoID(r.GetID())

	// send API call to create the pipeline
	p, err := database.FromContext(c).CreatePipeline(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create pipeline %s/%s: %w", r.GetFullName(), input.GetCommit(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"pipeline":    p.GetCommit(),
		"pipeline_id": p.GetID(),
	}).Info("pipeline created for repo")

	c.JSON(http.StatusCreated, p)
}
