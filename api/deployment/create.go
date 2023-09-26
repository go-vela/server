// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/deployments/{org}/{repo} deployments CreateDeployment
//
// Create a deployment for the configured backend
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the deployment
//     schema:
//       "$ref": "#/definitions/Deployment"
//   '400':
//     description: Unable to create the deployment
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the deployment
//     schema:
//       "$ref": "#/definitions/Error"

// CreateDeployment represents the API handler to
// create a deployment in the configured backend.
func CreateDeployment(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("creating new deployment for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Deployment)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new deployment for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in deployment object
	input.SetRepoID(r.GetID())
	input.SetUser(u.GetName())

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

	// send API call to create the deployment
	err = scm.FromContext(c).CreateDeployment(ctx, u, r, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create new deployment for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, input)
}
