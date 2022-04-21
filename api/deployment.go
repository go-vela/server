// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/router/middleware/org"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/deployments/{org}/{repo} deployment CreateDeployment
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

	// send API call to create the deployment
	err = scm.FromContext(c).CreateDeployment(u, r, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create new deployment for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, input)
}

// swagger:operation GET /api/v1/deployments/{org}/{repo} deployment GetDeployments
//
// Get a list of deployments for the configured backend
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
//     description: Successfully retrieved the list of deployments
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Deployment"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of deployments
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of deployments
//     schema:
//       "$ref": "#/definitions/Error"

// GetDeployments represents the API handler to capture
// a list of deployments from the configured backend.
func GetDeployments(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("reading deployments for repo %s", r.GetFullName())

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the total number of deployments for the repo
	t, err := scm.FromContext(c).GetDeploymentCount(u, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get deployment count for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the list of steps for the build
	d, err := scm.FromContext(c).GetDeploymentList(u, r, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get deployments for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	dWithBs := []*library.Deployment{}

	for _, deployment := range d {
		b, err := database.FromContext(c).GetDeploymentBuildList(*deployment.URL)
		if err != nil {
			retErr := fmt.Errorf("unable to get builds for deployment %d: %w", deployment.GetID(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		builds := []library.Build{}
		for _, build := range b {
			builds = append(builds, *build)
		}

		deployment.SetBuilds(builds)

		dWithBs = append(dWithBs, deployment)
	}

	// create pagination object
	pagination := Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, dWithBs)
}

// swagger:operation GET /api/v1/deployments/{org}/{repo}/{deployment} deployment GetDeployment
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
//   description: Name of the org
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
	deployment := c.Param("deployment")

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
	d, err := scm.FromContext(c).GetDeployment(u, r, int64(number))
	if err != nil {
		retErr := fmt.Errorf("unable to get deployment %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, d)
}
