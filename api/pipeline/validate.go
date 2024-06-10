// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/{pipeline}/validate pipelines ValidatePipeline
//
// Get, expand and validate a pipeline
//
// ---
// produces:
// - application/yaml
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
//     description: Successfully retrieved, expanded and validated the pipeline
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

// ValidatePipeline represents the API handler to capture,
// expand and validate a pipeline configuration.
func ValidatePipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
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
	}).Debugf("validating pipeline %s", entry)

	// ensure we use the expected pipeline type when compiling
	r.SetPipelineType(p.GetType())

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithCommit(p.GetCommit()).WithMetadata(m).WithRepo(r).WithUser(u)

	ruleData := prepareRuleData(c)

	// validate the pipeline
	pipeline, _, err := compiler.CompileLite(p.GetData(), ruleData, false)
	if err != nil {
		retErr := fmt.Errorf("unable to validate pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	writeOutput(c, pipeline)
}
