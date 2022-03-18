// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/pipelines/{org}/{repo}/{pipeline} pipelines UpdatePipeline
//
// update thes a pipeline in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: pipeline
//   description: Pipeline number to update
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing the pipeline to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Pipeline"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the pipeline
//     schema:
//       "$ref": "#/definitions/Pipeline"
//   '404':
//     description: Unable to update the pipeline
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the pipeline
//     schema:
//       "$ref": "#/definitions/Error"

// UpdatePipeline represents the API handler to update
// a pipeline for a repo in the configured backend.
func UpdatePipeline(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), p.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"pipeline": p.GetNumber(),
		"repo":     r.GetName(),
		"user":     u.GetName(),
	}).Infof("updating pipeline %s", entry)

	// capture body from API request
	input := new(library.Pipeline)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// check if the Commit field in the pipeline was provided
	if len(input.GetCommit()) > 0 {
		// update the Commit field
		p.SetCommit(input.GetCommit())
	}

	// check if the Flavor field in the pipeline was provided
	if len(input.GetFlavor()) > 0 {
		// update the Flavor field
		p.SetFlavor(input.GetFlavor())
	}

	// check if the Platform field in the pipeline was provided
	if len(input.GetPlatform()) > 0 {
		// update the Platform field
		p.SetPlatform(input.GetPlatform())
	}

	// check if the Ref field in the pipeline was provided
	if len(input.GetRef()) > 0 {
		// update the Ref field
		p.SetRef(input.GetRef())
	}

	// check if the Type field in the pipeline was provided
	if len(input.GetType()) > 0 {
		// update the Type field
		p.SetType(input.GetType())
	}

	// check if the Version field in the pipeline was provided
	if len(input.GetVersion()) > 0 {
		// update the Version field
		p.SetVersion(input.GetVersion())
	}

	// check if the ExternalSecrets field in the pipeline was provided
	if input.ExternalSecrets != nil {
		// update the ExternalSecrets field
		p.SetExternalSecrets(input.GetExternalSecrets())
	}

	// check if the InternalSecrets field in the pipeline was provided
	if input.InternalSecrets != nil {
		// update the InternalSecrets field
		p.SetInternalSecrets(input.GetInternalSecrets())
	}

	// check if the Services field in the pipeline was provided
	if input.Services != nil {
		// update the Services field
		p.SetServices(input.GetServices())
	}

	// check if the Stages field in the pipeline was provided
	if input.Stages != nil {
		// update the Stages field
		p.SetStages(input.GetStages())
	}

	// check if the Steps field in the pipeline was provided
	if input.Steps != nil {
		// update the Steps field
		p.SetSteps(input.GetSteps())
	}

	// check if the Templates field in the pipeline was provided
	if input.Templates != nil {
		// update the Templates field
		p.SetTemplates(input.GetTemplates())
	}

	// check if the Data field in the pipeline was provided
	if len(input.GetData()) > 0 {
		// update the data field
		p.SetData(input.GetData())
	}

	// send API call to update the pipeline
	err = database.FromContext(c).UpdatePipeline(p)
	if err != nil {
		retErr := fmt.Errorf("unable to update pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated pipeline
	p, err = database.FromContext(c).GetPipelineForRepo(p.GetNumber(), r)
	if err != nil {
		retErr := fmt.Errorf("unable to capture pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, p)
}
