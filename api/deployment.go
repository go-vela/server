// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/source"
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
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '201':
//     description: Successfully created the deployment
//     schema:
//       type: string
//   '400':
//     description: Successfully created the deployment
//     schema:
//       type: string
//   '500':
//     description: Successfully created the deployment
//     schema:
//       type: string

// CreateDeployment represents the API handler to
// create a deployment in the configured backend.
func CreateDeployment(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	logrus.Infof("Creating new deployment for %s", r.GetFullName())

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
	err = source.FromContext(c).CreateDeployment(u, r, input)
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
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the list of deployments
//     schema:
//       type: string

// GetDeployments represents the API handler to capture
// a list of deployments from the configured backend.
func GetDeployments(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	logrus.Infof("Reading deployments for %s", r.GetFullName())

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
	t, err := source.FromContext(c).GetDeploymentCount(u, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get deployment count for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the list of steps for the build
	d, err := source.FromContext(c).GetDeploymentList(u, r, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get deployments for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, d)
}

// swagger:operation GET /api/v1/deployments/{org}/{repo}/{deployment} deployment GetDeployment
//
// Get a deployment from the configured backend
//
// ---
// x-success_http_code: '501'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: deployment
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '502':
//     description: Successfully retrieved the deployment
//     schema:
//       type: string

// GetDeployment represents the API handler to
// capture a deployment from the configured backend.
func GetDeployment(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	deployment := c.Param("deployment")

	logrus.Infof("Reading deployment %s/%s", r.GetFullName(), deployment)

	number, err := strconv.Atoi(deployment)
	if err != nil {
		retErr := fmt.Errorf("invalid deployment parameter provided: %s", deployment)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the deployment
	d, err := source.FromContext(c).GetDeployment(u, r, int64(number))
	if err != nil {
		retErr := fmt.Errorf("unable to get deployment %s/%d: %w", r.GetFullName(), number, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, d)
}
