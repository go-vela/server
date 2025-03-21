// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/deployments/{org}/{repo}/{deployment} deployments GetDeployment
//
// Get a deployment
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
//   name: deployment
//   description: Number of the deployment
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the deployment
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

// GetDeployment represents the API handler to get a deployment.
func GetDeployment(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	deployment := util.PathParameter(c, "deployment")
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), deployment)

	l.Debugf("reading deployment %s", entry)

	number, err := strconv.Atoi(deployment)
	if err != nil {
		retErr := fmt.Errorf("invalid deployment parameter provided: %s", deployment)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to database to capture the deployment
	d, err := database.FromContext(c).GetDeploymentForRepo(ctx, r, int64(number))
	if err != nil {
		// send API call to SCM to capture the deployment
		d, err = scm.FromContext(c).GetDeployment(ctx, u, r, int64(number))
		if err != nil {
			retErr := fmt.Errorf("unable to get deployment %s: %w", entry, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	if d == nil {
		retErr := fmt.Errorf("unable to get deployment: %s", deployment)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	c.JSON(http.StatusOK, d)
}
