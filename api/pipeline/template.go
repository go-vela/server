// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/registry/github"
	"github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/pipelines/{org}/{repo}/{pipeline}/templates pipelines GetTemplates
//
// Get pipeline templates
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
//     description: Successfully retrieved the map of pipeline templates
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Template"
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

// GetTemplates represents the API handler to capture a
// map of templates utilized by a pipeline configuration.
func GetTemplates(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	o := org.Retrieve(c)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), p.GetCommit())

	l.Debugf("reading templates from pipeline %s", entry)

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithCommit(p.GetCommit()).WithMetadata(m).WithRepo(r).WithUser(u)

	// parse the pipeline configuration
	pipeline, _, err := compiler.Parse(p.GetData(), p.GetType(), new(yaml.Template))
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, fmt.Errorf("unable to parse pipeline %s: %w", entry, err))

		return
	}

	baseErr := fmt.Sprintf("unable to set template links for %s", entry)

	templates := make(map[string]*types.Template)
	for name, template := range pipeline.Templates.Map() {
		templates[name] = template.ToAPI()

		// create a compiler registry client for parsing (no address or token needed for Parse)
		registry, err := github.New(ctx, "", "")
		if err != nil {
			util.HandleError(c, http.StatusBadRequest, fmt.Errorf("%s: unable to create compiler github client: %w", baseErr, err))

			return
		}

		// capture source path to template
		source := template.Source

		// if type is file, compose a source string so the template can be found
		if strings.EqualFold(template.Type, "file") {
			source = fmt.Sprintf("%s%s/%s/%s@%s", registry.URL, o, r.GetName(), source, p.GetCommit())
		}

		// parse the source for the template using the compiler registry client
		src, err := registry.Parse(source)
		if err != nil {
			util.HandleError(c, http.StatusBadRequest, fmt.Errorf("%s: unable to parse source for %s: %w", baseErr, template.Source, err))

			return
		}

		// retrieve link to template file from github
		link, err := scm.FromContext(c).GetHTMLURL(ctx, r.GetOwner(), src.Org, src.Repo, src.Name, src.Ref)
		if err != nil {
			util.HandleError(c, http.StatusBadRequest, fmt.Errorf("%s: unable to get html url for %s/%s/%s/@%s: %w", baseErr, src.Org, src.Repo, src.Name, src.Ref, err))

			return
		}

		// set link to file for template
		templates[name].SetLink(link)
	}

	writeOutput(c, templates)
}
