package api

import (
	"fmt"
	"net/http"

	"github.com/go-vela/compiler/compiler"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types"
	"github.com/go-vela/types/yaml"

	"github.com/gin-gonic/gin"
)

func Compile(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), "master")
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// parse and compile the pipeline configuration file
	p, err := compiler.FromContext(c).
		// WithBuild(input).
		// WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		// Compile(config)
		// parse the object into a yaml configuration
		Parse(config)

	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create map of templates for easy lookup
	tmpls := mapFromTemplates(p.Templates)

	if len(p.Stages) > 0 {
		// inject the templates into the stages
		p.Stages, err = compiler.FromContext(c).
			// WithBuild(input).
			// WithFiles(files).
			WithMetadata(m).
			WithRepo(r).
			WithUser(u).
			// Compile(config)
			// parse the object into a yaml configuration
			ExpandStages(p.Stages, tmpls)
		if err != nil {
			retErr := fmt.Errorf("unable to compile pipeline configuration for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
		c.YAML(http.StatusEarlyHints, p)
		return
	}

	// inject the templates into the stages
	p.Steps, err = compiler.FromContext(c).
		// WithBuild(input).
		// WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		// Compile(config)
		// parse the object into a yaml configuration
		ExpandSteps(p.Steps, tmpls)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.YAML(http.StatusCreated, p)
}

// helper function that creates a map of templates from a yaml configuration.
func mapFromTemplates(templates []*yaml.Template) map[string]*yaml.Template {
	m := make(map[string]*yaml.Template)

	for _, tmpl := range templates {
		m[tmpl.Name] = tmpl
	}

	return m
}
