// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/deployments/{org}/{repo}/config deployments GetDeploymentConfig
//
// Get a deployment config
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
//   name: ref
//   description: Ref to target for the deployment config
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the deployment config
//     schema:
//       "$ref": "#/definitions/Deployment"
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

// GetDeploymentConfig represents the API handler to get a deployment config at a given ref.
func GetDeploymentConfig(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	ctx := c.Request.Context()

	// capture ref from parameters - use default branch if not provided
	ref := util.QueryParameter(c, "ref", r.GetBranch())

	entry := fmt.Sprintf("%s@%s", r.GetFullName(), ref)

	l.Debugf("reading deployment config %s", entry)

	var config []byte

	// check if the pipeline exists in the database
	p, err := database.FromContext(c).GetPipelineForRepo(ctx, ref, r)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Debugf("pipeline %s not found in database, fetching from scm", entry)

			config, err = scm.FromContext(c).ConfigBackoff(ctx, u, r, ref)
			if err != nil {
				retErr := fmt.Errorf("unable to get pipeline configuration for %s: %w", entry, err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}
		} else {
			// some other error
			retErr := fmt.Errorf("unable to get pipeline for %s: %w", entry, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	} else {
		l.Debugf("pipeline %s found in database", entry)

		config = p.GetData()
	}

	// set up compiler
	compiler := compiler.FromContext(c).Duplicate().WithCommit(ref).WithRepo(r).WithUser(u)

	// compile the pipeline
	pipeline, _, err := compiler.CompileLite(ctx, config, nil, true)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	c.JSON(http.StatusOK, pipeline.Deployment)
}