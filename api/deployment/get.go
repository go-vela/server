// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/deployments/{org}/{repo}/{deployment} deployments GetDeployment
//
// Get a deployment from the configured backend
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
//     description: Unable to retrieve the deployment
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the deployment
//     schema:
//       "$ref": "#/definitions/Error"

// GetDeployment represents the API handler to
// capture a deployment from the configured backend.
func GetDeployment(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	deployment := util.PathParameter(c, "deployment")

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), deployment)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("reading deployment %s", entry)

	number, err := strconv.Atoi(deployment)
	if err != nil {
		retErr := fmt.Errorf("invalid deployment parameter provided: %s", deployment)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the deployment
	d, err := database.FromContext(c).GetDeployment(c, int64(number))
	if err != nil {
		retErr := fmt.Errorf("unable to get deployment %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, d)
}
