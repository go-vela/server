// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/pipelines/{org}/{repo} pipelines ListPipelines
//
// Get all pipelines for a repository
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
// - in: query
//   name: page
//   description: The page of results to retrieve
//   type: integer
//   default: 1
// - in: query
//   name: per_page
//   description: How many results per page to return
//   type: integer
//   maximum: 100
//   default: 10
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the pipelines
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Pipeline"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: See https://tools.ietf.org/html/rfc5988
//         type: string
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

// ListPipelines represents the API handler to get a list
// of pipelines for a repository.
func ListPipelines(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("listing pipelines for repo %s", r.GetFullName())

	pagination, err := api.ParsePagination(c)
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, err)
		return
	}

	p, err := database.FromContext(c).ListPipelinesForRepo(ctx, r, pagination.Page, pagination.PerPage)
	if err != nil {
		retErr := fmt.Errorf("unable to list pipelines for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// set pagination results
	pagination.Results = len(p)
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, p)
}
