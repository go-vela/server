// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"fmt"
	"net/http"
	"strconv"

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

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/{pipeline}/validate pipelines ValidatePipeline
//
// Get, expand and validate a pipeline from the configured backend
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

// ValidatePipeline represents the API handler to capture,
// expand and validate a pipeline configuration.
func ValidatePipeline(c *gin.Context) {
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
	}).Infof("validating pipeline %s", entry)

	// ensure we use the expected pipeline type when compiling
	r.SetPipelineType(p.GetType())

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithMetadata(m).WithRepo(r).WithUser(u)

	// capture optional template query parameter
	template, err := strconv.ParseBool(c.DefaultQuery("template", "true"))
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, fmt.Errorf("unable to parse template query parameter for %s: %w", entry, err))

		return
	}

	// validate the pipeline
	pipeline, _, err := compiler.CompileLite(p.GetData(), template, false, nil)
	if err != nil {
		retErr := fmt.Errorf("unable to validate pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	writeOutput(c, pipeline)
}
