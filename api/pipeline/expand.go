// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with compile
package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/{pipeline}/expand pipelines ExpandPipeline
//
// Get and expand a pipeline from the configured backend
//
// ---
// produces:
// - application/x-yaml
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: pipeline
//   description: Commit SHA for pipeline to retrieve
//   required: true
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
//   type: string
//   default: yaml
//   enum:
//   - json
//   - yaml
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved and expanded the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
//   '400':
//     description: Unable to expand the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"

// ExpandPipeline represents the API handler to capture and
// expand a pipeline configuration.
func ExpandPipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	o := org.Retrieve(c)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), p.GetCommit())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"pipeline": p.GetCommit(),
		"repo":     r.GetName(),
		"user":     u.GetName(),
	}).Infof("expanding templates for pipeline %s", entry)

	// ensure we use the expected pipeline type when compiling
	r.SetPipelineType(p.GetType())

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithCommit(p.GetCommit()).WithMetadata(m).WithRepo(r).WithUser(u)

	// expand the templates in the pipeline
	pipeline, _, err := compiler.CompileLite(p.GetData(), false)
	if err != nil {
		retErr := fmt.Errorf("unable to expand pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	writeOutput(c, pipeline)
}
