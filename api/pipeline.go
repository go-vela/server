package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/compiler/compiler"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types"
	"github.com/go-vela/types/yaml"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /api/v1/pipelines/{org}/{repo} pipeline GetPipeline
//
//
// ---
// x-success_http_code: '200'
// produces:
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   required: true
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format (default: yaml)
//   required: false
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/Pipeline"
//   '400':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       type: string
//   '404':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       type: string

// GetPipeline represents the API handler to capture a
// pipeline configuration for a repo from the the source provider.
func GetPipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// capture query parameters
	output := c.DefaultQuery("output", "yaml")
	ref := c.DefaultQuery("ref", r.GetBranch())

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)

	// parse the pipeline configuration file
	p, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case "json":
		c.JSON(http.StatusOK, p)
	case "yaml":
		fallthrough
	default:
		c.YAML(http.StatusOK, p)
	}
}

// swagger:operation GET /api/v1/pipelines/{org}/{repo}/templates templates GetTemplates
//
//
// ---
// x-success_http_code: '200'
// produces:
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   required: true
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format (default: yaml)
//   required: false
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the map of pipeline templates
//     type: json
//     schema:
//       "$ref": "#/definitions/Template"
//   '400':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       type: string
//   '404':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       type: string

// GetTemplates represents the API handler to capture a
// map of templates utilized by a pipeline configuration.
func GetTemplates(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// capture query parameters
	output := c.DefaultQuery("output", "yaml")
	ref := c.DefaultQuery("ref", r.GetBranch())

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)

	// parse the pipeline configuration file
	p, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create map of templates for response body
	t := mapFromTemplates(p.Templates)

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case "json":
		c.JSON(http.StatusOK, t)
	case "yaml":
		fallthrough
	default:
		c.YAML(http.StatusOK, t)
	}
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/expand pipeline ExpandPipeline
//
//
// ---
// x-success_http_code: '200'
// produces:
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   required: true
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format (default: yaml)
//   required: false
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved and expanded the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/Pipeline"
//   '400':
//     description: Unable to expand the pipeline configuration
//     schema:
//       type: string
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       type: string
//   '500':
//     description: Unable to expand the pipeline configuration
//     schema:
//       type: string

// ExpandPipeline represents the API handler to capture and
// expand a pipeline configuration.
func ExpandPipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// capture query parameters
	output := c.DefaultQuery("output", "yaml")
	ref := c.DefaultQuery("ref", r.GetBranch())

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)

	// parse the pipeline configuration file
	p, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create map of templates for easy lookup
	t := mapFromTemplates(p.Templates)

	// check if the pipeline contains stages
	if len(p.Stages) > 0 {
		// inject the templates into the stages
		p.Stages, err = comp.ExpandStages(p.Stages, t)
		if err != nil {
			retErr := fmt.Errorf("unable to expand stages in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	} else {
		// inject the templates into the steps
		p.Steps, err = comp.ExpandSteps(p.Steps, t)
		if err != nil {
			retErr := fmt.Errorf("unable to expand steps in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case "json":
		c.JSON(http.StatusOK, p)
	case "yaml":
		fallthrough
	default:
		c.YAML(http.StatusOK, p)
	}
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/validate pipeline ValidatePipeline
//
//
// ---
// x-success_http_code: '200'
// produces:
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   required: true
//   type: string
// - in: query
//   name: template
//   description: Template boolean for expanding templates (default: true)
//   required: false
//   type: bool
// - in: query
//   name: output
//   description: Output string for specifying output format (default: yaml)
//   required: false
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved, expanded and validated the pipeline
//     schema:
//       type: string
//   '400':
//     description: Unable to validate the pipeline configuration
//     schema:
//       type: string
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       type: string
//   '500':
//     description: Unable to validate the pipeline configuration
//     schema:
//       type: string

// ValidatePipeline represents the API handler to capture, expand and
// validate a pipeline configuration.
func ValidatePipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// capture query parameters
	ref := c.DefaultQuery("ref", r.GetBranch())

	template := c.DefaultQuery("template", "true")

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)

	// parse the pipeline configuration file
	p, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// check optional template query parameter
	if strings.ToLower(template) == "true" {
		// create map of templates for easy lookup
		t := mapFromTemplates(p.Templates)

		// check if the pipeline contains stages
		if len(p.Stages) > 0 {
			// inject the templates into the stages
			p.Stages, err = comp.ExpandStages(p.Stages, t)
			if err != nil {
				retErr := fmt.Errorf("unable to expand stages in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}
		} else { // inject the templates into the stages
			p.Steps, err = comp.ExpandSteps(p.Steps, t)
			if err != nil {
				retErr := fmt.Errorf("unable to expand steps in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}
		}
	}

	// validate the yaml configuration
	err = comp.Validate(p)
	if err != nil {
		retErr := fmt.Errorf("unable to validate pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// pipeline is valid, respond to user
	c.JSON(http.StatusOK, "pipeline is valid")
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/compile pipeline CompilePipeline
//
//
// ---
// x-success_http_code: '200'
// produces:
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   required: true
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format (default: yaml)
//   required: false
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved and compiled the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/Pipeline"
//   '400':
//     description: Unable to validate the pipeline configuration
//     schema:
//       type: string
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       type: string
//   '500':
//     description: Unable to validate the pipeline configuration
//     schema:
//       type: string

// CompilePipeline represents the API handler to capture, expand and
// compile a pipeline configuration.
func CompilePipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// capture query parameters
	output := c.DefaultQuery("output", "yaml")
	ref := c.DefaultQuery("ref", r.GetBranch())

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)

	// parse the pipeline configuration file
	p, err := comp.Compile(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case "json":
		c.JSON(http.StatusOK, p)
	case "yaml":
		fallthrough
	default:
		c.YAML(http.StatusOK, p)
	}
}

// helper function that creates a map of templates from a yaml configuration.
func mapFromTemplates(templates []*yaml.Template) map[string]*yaml.Template {
	m := make(map[string]*yaml.Template)

	for _, tmpl := range templates {
		m[tmpl.Name] = tmpl
	}

	return m
}
