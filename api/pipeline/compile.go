// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with expand
package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/internal"
	pMiddleware "github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/pipeline"
)

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/{pipeline}/compile pipelines CompilePipeline
//
// Get, expand and compile a pipeline
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
//     description: Successfully retrieved and compiled the pipeline
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
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

// CompilePipeline represents the API handler to capture,
// expand and compile a pipeline configuration.
func CompilePipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	p := pMiddleware.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), p.GetCommit())

	l.Debugf("compiling pipeline %s", entry)

	// ensure we use the expected pipeline type when compiling
	r.SetPipelineType(p.GetType())

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithCommit(p.GetCommit()).WithMetadata(m).WithRepo(r).WithUser(u)

	ruleData := prepareRuleData(c)

	// compile the pipeline
	pipeline, _, err := compiler.CompileLite(ctx, p.GetData(), ruleData, true)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	writeOutput(c, pipeline)
}

// prepareRuleData is a helper function to prepare the rule data from the query parameters.
func prepareRuleData(c *gin.Context) *pipeline.RuleData {
	// capture the branch name parameter
	branch := c.Query("branch")
	// capture the comment parameter
	comment := c.Query("comment")
	// capture the event type parameter
	event := c.Query("event")
	// capture the repo parameter
	ruleDataRepo := c.Query("repo")
	// capture the status type parameter
	status := c.Query("status")
	// capture the tag parameter
	tag := c.Query("tag")
	// capture the target parameter
	target := c.Query("target")
	// capture the path parameter
	pathSet := c.QueryArray("path")

	// if any ruledata query params were provided, create ruledata struct
	if len(branch) > 0 ||
		len(comment) > 0 ||
		len(event) > 0 ||
		len(pathSet) > 0 ||
		len(ruleDataRepo) > 0 ||
		len(status) > 0 ||
		len(tag) > 0 ||
		len(target) > 0 {
		return &pipeline.RuleData{
			Branch:  branch,
			Comment: comment,
			Event:   event,
			Path:    pathSet,
			Repo:    ruleDataRepo,
			Status:  status,
			Tag:     tag,
			Target:  target,
		}
	}

	return nil
}
