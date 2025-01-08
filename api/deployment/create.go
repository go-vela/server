// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/deployments/{org}/{repo} deployments CreateDeployment
//
// Create a deployment
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the deployment
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

// CreateDeployment represents the API handler to
// create a deployment.
func CreateDeployment(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("creating new deployment for repo %s", r.GetFullName())

	// capture body from API request
	input := new(types.Deployment)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new deployment for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in deployment object
	input.SetRepo(r)
	input.SetCreatedBy(u.GetName())
	input.SetCreatedAt(time.Now().Unix())

	if len(input.GetDescription()) == 0 {
		input.SetDescription("Deployment request from Vela")
	}

	if len(input.GetTask()) == 0 {
		input.SetTask("deploy:vela")
	}

	// if ref is not provided, use repo default branch
	if len(input.GetRef()) == 0 {
		input.SetRef(fmt.Sprintf("refs/heads/%s", r.GetBranch()))
	}

	deployConfigYAML, err := getDeploymentConfig(c, l, u, r, input.GetRef())
	if err != nil {
		retErr := fmt.Errorf("unable to get deployment config for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	deployConfig := deployConfigYAML.ToPipeline()

	if !deployConfig.Empty() {
		err := deployConfig.Validate(input.GetTarget(), input.GetPayload())

		if err != nil {
			retErr := fmt.Errorf("unable to validate deployment config for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	// send API call to create the deployment
	err = scm.FromContext(c).CreateDeployment(ctx, u, r, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create new deployment for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to create the deployment
	d, err := database.FromContext(c).CreateDeployment(c, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create new deployment for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"deployment_id": d.GetID(),
	}).Info("deployment created")

	c.JSON(http.StatusCreated, d)
}
