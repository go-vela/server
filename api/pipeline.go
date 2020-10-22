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

func ExpandPipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// capture query parameters
	ref := c.DefaultQuery("ref", r.GetBranch())

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// parse the pipeline configuration file
	p, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

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

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		c.YAML(http.StatusOK, p)
		return
	}

	// inject the templates into the stages
	p.Steps, err = comp.ExpandSteps(p.Steps, t)
	if err != nil {
		retErr := fmt.Errorf("unable to expand steps in pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.YAML(http.StatusOK, p)
}

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

	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), ref)
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// parse the pipeline configuration file
	p, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

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

func GetPipeline(c *gin.Context) {
	// capture middleware values
	// m := c.MustGet("metadata").(*types.Metadata)
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

	m := c.MustGet("metadata").(*types.Metadata)
	// create the compiler with extra information embedded into it
	comp := compiler.FromContext(c).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u)
	// parse the pipeline configuration file
	p, err := comp.Parse(config)
	if err != nil {
		retErr := fmt.Errorf("unable to parse pipeline configuration for %s@%s: %w", r.GetFullName(), ref, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

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

// helper function that creates a map of templates from a yaml configuration.
func mapFromTemplates(templates []*yaml.Template) map[string]*yaml.Template {
	m := make(map[string]*yaml.Template)

	for _, tmpl := range templates {
		m[tmpl.Name] = tmpl
	}

	return m
}
