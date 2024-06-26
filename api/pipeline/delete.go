// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/pipelines/{org}/{repo}/{pipeline} pipelines DeletePipeline
//
// Delete a pipeline
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
//   description: Commit SHA for pipeline to delete
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the pipeline
//     schema:
//       type: string
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

// DeletePipeline represents the API handler to remove a pipeline for a repository.
func DeletePipeline(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), p.GetCommit())

	l.Debugf("deleting pipeline %s", entry)

	// send API call to remove the build
	err := database.FromContext(c).DeletePipeline(ctx, p)
	if err != nil {
		retErr := fmt.Errorf("unable to delete pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("pipeline %s deleted", entry))
}
