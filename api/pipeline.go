// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/registry/github"

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
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
//   '400':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       "$ref": "#/definitions/Error"

// GetPipeline represents the API handler to capture a
// pipeline configuration for a repo from the the source provider.
func GetPipeline(ctx *gin.Context) {
	pipeline, _, err := getUnprocessedPipeline(ctx)
	if err != nil {
		util.HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	writeOutput(ctx, pipeline)
}

// swagger:operation GET /api/v1/pipelines/{org}/{repo}/templates pipelines GetTemplates
//
// Get a map of templates utilized by a pipeline configuration from the source provider
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
//     schema:
//       "$ref": "#/definitions/Template"
//   '400':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       "$ref": "#/definitions/Error"

// GetTemplates represents the API handler to capture a
// map of templates utilized by a pipeline configuration.
func GetTemplates(ctx *gin.Context) {
	pipeline, _, err := getUnprocessedPipeline(ctx)
	if err != nil {
		util.HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	// create map of templates for response body
	templates, err := getTemplateLinks(ctx, pipeline.Templates)
	if err != nil {
		retErr := fmt.Errorf("unable to set template links for %s: %w", repoName(ctx), err)
		util.HandleError(ctx, http.StatusBadRequest, retErr)
		return
	}

	writeOutput(ctx, templates)
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/expand pipelines ExpandPipeline
//
// Get and expand a pipeline configuration from the source provider
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
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"

// ExpandPipeline represents the API handler to capture and
// expand a pipeline configuration.
func ExpandPipeline(ctx *gin.Context) {
	pipeline, comp, err := getUnprocessedPipeline(ctx)
	if err != nil {
		util.HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := expandPipeline(ctx, pipeline, comp, false); err != nil {
		util.HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	writeOutput(ctx, pipeline)
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/validate pipelines ValidatePipeline
//
// Get, expand and validate a pipeline configuration from the source provider
//
// ---
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
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"

// ValidatePipeline represents the API handler to capture, expand and
// validate a pipeline configuration.
func ValidatePipeline(ctx *gin.Context) {
	pipeline, comp, err := getUnprocessedPipeline(ctx)
	if err != nil {
		util.HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	// check optional template query parameter
	if ok, _ := strconv.ParseBool(ctx.DefaultQuery("template", "true")); ok {
		if err := expandPipeline(ctx, pipeline, comp, false); err != nil {
			util.HandleError(ctx, http.StatusBadRequest, err)
			return
		}
	}

	// validate the yaml configuration
	if err = comp.Validate(pipeline); err != nil {
		retErr := fmt.Errorf("unable to validate pipeline configuration for %s: %w", repoName(ctx), err)
		util.HandleError(ctx, http.StatusBadRequest, retErr)
		return
	}

	writeOutput(ctx, pipeline)
}

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/compile pipelines CompilePipeline
//
// Get, expand and compile a pipeline configuration from the source provider
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
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
//   '400':
//     description: Unable to validate the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"

// CompilePipeline represents the API handler to capture,
// expand and compile a pipeline configuration.
//
func CompilePipeline(ctx *gin.Context) {
	pipeline, comp, err := getUnprocessedPipeline(ctx)
	if err != nil {
		util.HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := expandPipeline(ctx, pipeline, comp, true); err != nil {
		util.HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	// validate the yaml configuration
	if err = comp.Validate(pipeline); err != nil {
		retErr := fmt.Errorf("unable to validate pipeline configuration for %s: %w", repoName(ctx), err)
		util.HandleError(ctx, http.StatusBadRequest, retErr)
		return
	}

	writeOutput(ctx, pipeline)
}

// getUnprocessedPipeline retrieves the unprocessed pipeline from a given context.
func getUnprocessedPipeline(ctx *gin.Context) (*yaml.Build, compiler.Engine, error) {
	// capture middleware values
	meta := ctx.MustGet("metadata").(*types.Metadata)
	repo := repo.Retrieve(ctx)

	// capture query parameters
	ref := ctx.DefaultQuery("ref", repo.GetBranch())

	// send API call to capture the repo owner
	user, err := database.FromContext(ctx).GetUser(repo.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", repo.GetFullName(), err)
		return nil, nil, retErr
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(ctx).ConfigBackoff(user, repo, ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s: %w", repoName(ctx), err)
		return nil, nil, retErr
	}

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(ctx).
		WithMetadata(meta).
		WithRepo(repo).
		WithUser(user)

	pipeline, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s: %w", repoName(ctx), err)
		return nil, nil, retErr
	}

	return pipeline, comp, nil
}

// getTemplateLinks helper function that retrieves source provider links
// for a list of templates and returns a map of library templates.
//
// nolint: lll // ignore long line length due to variable names
func getTemplateLinks(ctx *gin.Context, templates yaml.TemplateSlice) (map[string]*library.Template, error) {
	r := repo.Retrieve(ctx)
	u, err := database.FromContext(ctx).GetUser(r.GetUserID())
	if err != nil {
		return nil, err
	}

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
		src, err := cl.Parse(tmpl.GetSource())
		if err != nil {
			retErr := fmt.Errorf("unable to parse source for %s: %w", tmpl.GetSource(), err)
			return nil, retErr
		}

		// retrieve link to template file from github
		link, err := source.FromContext(ctx).GetHTMLURL(u, src.Org, src.Repo, src.Name, src.Ref)
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

// repoName takes the given context and returns a string friendly
// representation with the format of 'repository@reference'.
func repoName(ctx *gin.Context) string {
	repo := repo.Retrieve(ctx)
	ref := ctx.DefaultQuery("ref", repo.GetBranch())

	return fmt.Sprintf("%s@%s", repo.GetFullName(), ref)
}

// writeOutput returns writes output to the request based on the preferred
// output as defined in the request's 'output' query defaulting to YAML.
func writeOutput(ctx *gin.Context, pipeline interface{}) {
	output := ctx.DefaultQuery("output", outputYAML)

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case outputJSON:
		ctx.JSON(http.StatusOK, pipeline)
	case outputYAML:
		fallthrough
	default:
		ctx.YAML(http.StatusOK, pipeline)
	}
}

// expandPipeline uses a given pipeline and compiler to expand stages and steps
// in the pipeline along with optionally substituting the environmental variables.
//
// nolint: lll // ignore long line length due to variable names
func expandPipeline(ctx *gin.Context, pipeline *yaml.Build, comp compiler.Engine, substituteEnv bool) error {
	// create map of templates for easy lookup
	templates := pipeline.Templates.Map()

	var err error

	if len(pipeline.Stages) > 0 {
		// inject the templates into the stages
		pipeline.Stages, pipeline.Secrets, pipeline.Services, pipeline.Environment, err = comp.ExpandStages(pipeline, templates)
		if err != nil {
			return fmt.Errorf("unable to expand stages in pipeline configuration for %s: %w", repoName(ctx), err)
		}

		if substituteEnv {
			// inject the substituted environment variables into the stages
			pipeline.Stages, err = comp.SubstituteStages(pipeline.Stages)
			if err != nil {
				// nolint: lll // ignore long line length due to error message
				return fmt.Errorf("unable to substitute stages in pipeline configuration for %s: %w", repoName(ctx), err)
			}
		}
	} else {
		// inject the templates into the steps
		pipeline.Steps, pipeline.Secrets, pipeline.Services, pipeline.Environment, err = comp.ExpandSteps(pipeline, templates)
		if err != nil {
			return fmt.Errorf("unable to expand steps in pipeline configuration for %s: %w", repoName(ctx), err)
		}

		if substituteEnv {
			// inject the substituted environment variables into the steps
			pipeline.Steps, err = comp.SubstituteSteps(pipeline.Steps)
			if err != nil {
				// nolint: lll // ignore long line length due to error message
				return fmt.Errorf("unable to substitute steps in pipeline configuration for %s: %w", repoName(ctx), err)
			}
		}
	}

	return nil
}
