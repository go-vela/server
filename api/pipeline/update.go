// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/pipelines/{org}/{repo}/{pipeline} pipelines UpdatePipeline
//
// Update a pipeline
//
// ---
// produces:
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
//   description: Commit SHA for pipeline to update
//   required: true
//   type: string
// - in: body
//   name: body
//   description: The pipeline object with the fields to be updated
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

// UpdatePipeline represents the API handler to update
// a pipeline for a repo.
func UpdatePipeline(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), p.GetCommit())

	l.Debugf("updating pipeline %s", entry)

	// capture body from API request
	input := new(types.Pipeline)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
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
	p, err = database.FromContext(c).UpdatePipeline(ctx, p)
	if err != nil {
		retErr := fmt.Errorf("unable to update pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, p)
}
