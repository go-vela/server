// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/compiler/compiler"
	"github.com/go-vela/compiler/registry/github"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/yaml"

	"github.com/gin-gonic/gin"
)

const (
	outputJSON = "json"
	outputYAML = "yaml"
)

// swagger:operation GET /api/v1/pipelines/{org}/{repo} pipelines GetPipeline
//
// Get a pipeline configuration from the source provider
//
// ---
// x-success_http_code: '200'
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
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
	output := c.DefaultQuery("output", outputYAML)
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
	p, _, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case outputJSON:
		c.JSON(http.StatusOK, p)
	case outputYAML:
		fallthrough
	default:
		c.YAML(http.StatusOK, p)
	}
}

// swagger:operation GET /api/v1/pipelines/{org}/{repo}/templates pipelines GetTemplates
//
// Get a map of templates utilized by a pipeline configuration from the source provider
//
// ---
// x-success_http_code: '200'
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
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
	output := c.DefaultQuery("output", outputYAML)
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
	p, _, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create map of templates for response body
	t, err := setTemplateLinks(c, u, r, p.Templates)
	if err != nil {
		retErr := fmt.Errorf("unable to set template links for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case outputJSON:
		c.JSON(http.StatusOK, t)
	case outputYAML:
		fallthrough
	default:
		c.YAML(http.StatusOK, t)
	}
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/expand pipelines ExpandPipeline
//
// Get and expand a pipeline configuration from the source provider
//
// ---
// x-success_http_code: '200'
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
//   type: string
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
	output := c.DefaultQuery("output", outputYAML)
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
	p, _, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create map of templates for easy lookup
	t := p.Templates.Map()

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
	case outputJSON:
		c.JSON(http.StatusOK, p)
	case outputYAML:
		fallthrough
	default:
		c.YAML(http.StatusOK, p)
	}
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/validate pipelines ValidatePipeline
//
// Get, expand and validate a pipeline configuration from the source provider
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
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
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
	p, raw, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// check optional template query parameter
	if strings.ToLower(template) == "true" {
		// create map of templates for easy lookup
		t := p.Templates.Map()

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
	err = p.Validate(raw)
	if err != nil {
		retErr := fmt.Errorf("unable to validate pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// pipeline is valid, respond to user
	c.JSON(http.StatusOK, "pipeline is valid")
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/compile pipelines CompilePipeline
//
// Get, expand and compile a pipeline configuration from the source provider
//
// ---
// x-success_http_code: '200'
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
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved and compiled the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
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

// CompilePipeline represents the API handler to capture,
// expand and compile a pipeline configuration.
func CompilePipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// capture query parameters
	output := c.DefaultQuery("output", outputYAML)
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
	p, raw, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create map of templates for easy lookup
	t := p.Templates.Map()

	// check if the pipeline contains stages
	if len(p.Stages) > 0 {
		// inject the templates into the stages
		p.Stages, err = comp.ExpandStages(p.Stages, t)
		if err != nil {
			retErr := fmt.Errorf("unable to expand stages in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// inject the substituted environment variables into the stages
		p.Stages, err = comp.SubstituteStages(p.Stages)
		if err != nil {
			retErr := fmt.Errorf("unable to substitute stages in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

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

		// inject the substituted environment variables into the steps
		p.Steps, err = comp.SubstituteSteps(p.Steps)
		if err != nil {
			retErr := fmt.Errorf("unable to substitute steps in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	// validate the yaml configuration
	err = p.Validate(raw)
	if err != nil {
		retErr := fmt.Errorf("unable to validate pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case outputJSON:
		c.JSON(http.StatusOK, p)
	case outputYAML:
		fallthrough
	default:
		c.YAML(http.StatusOK, p)
	}
}

// setTemplateLinks helper function that retrieves source provider links for a list of templates and returns a map of library templates.
func setTemplateLinks(c *gin.Context, u *library.User, r *library.Repo, templates yaml.TemplateSlice) (map[string]*library.Template, error) {
	m := make(map[string]*library.Template)
	for _, t := range templates {

		// convert to library type
		tmpl := t.ToLibrary()

		// create a new compiler github client for parsing,
		// no address or token needed for Parse
		cl, err := github.New("", "")
		if err != nil {
			retErr := fmt.Errorf("unable to create compiler github client: %w", err)

			return nil, retErr
		}

		// parse template source
		src, err := cl.Parse(tmpl.GetSource(), r.GetBranch())
		if err != nil {
			retErr := fmt.Errorf("unable to parse source for %s: %w", tmpl.GetSource(), err)

			return nil, retErr
		}

		// retrieve link to template file from github
		link, err := source.FromContext(c).GetHTMLURL(u, src.Org, src.Repo, src.Name, src.Ref)
		if err != nil {
			retErr := fmt.Errorf("unable to get html url for %s/%s/%s/@%s: %w", src.Org, src.Repo, src.Name, src.Ref, err)

			return nil, retErr
		}

		// set link to template file
		tmpl.SetLink(link)

		m[tmpl.GetName()] = tmpl
	}

	return m, nil
}
