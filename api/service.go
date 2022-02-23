// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/services services CreateService
//
// Create a service for a build in the configured backend
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
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing the service to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Service"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the service
//     schema:
//       "$ref": "#/definitions/Service"
//   '400':
//     description: Unable to create the service
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the service
//     schema:
//       "$ref": "#/definitions/Error"

// CreateService represents the API handler to create
// a service for a build in the configured backend.
//
// nolint: dupl // ignore similar code with step
func CreateService(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("creating new service for build %s", entry)

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new service for build %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in service object
	input.SetRepoID(r.GetID())
	input.SetBuildID(b.GetID())

	if len(input.GetStatus()) == 0 {
		input.SetStatus(constants.StatusPending)
	}

	if input.GetCreated() == 0 {
		input.SetCreated(time.Now().UTC().Unix())
	}

	// send API call to create the service
	err = database.FromContext(c).CreateService(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create service for build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created service
	s, _ := database.FromContext(c).GetService(input.GetNumber(), b)

	c.JSON(http.StatusCreated, s)
}

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/services services GetServices
//
// Get a list of all services for a build in the configured backend
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
//   name: build
//   description: Build number
//   required: true
//   type: integer
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
//     description: Successfully retrieved the list of services
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Service"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of services
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of services
//     schema:
//       "$ref": "#/definitions/Error"

// GetServices represents the API handler to capture a list
// of services for a build from the configured backend.
func GetServices(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("reading services for build %s", entry)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for build %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for build %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	//
	// nolint: gomnd // ignore magic number
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the total number of services for the build
	t, err := database.FromContext(c).GetBuildServiceCount(b)
	if err != nil {
		retErr := fmt.Errorf("unable to get services count for build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the list of services for the build
	s, err := database.FromContext(c).GetBuildServiceList(b, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get services for build %s: %w", entry, err)

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

	c.JSON(http.StatusOK, s)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/services/{service} services GetService
//
// Get a service for a build in the configured backend
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
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: service
//   description: Name of the service
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the service
//     schema:
//       "$ref": "#/definitions/Service"
//   '400':
//     description: Unable to retrieve the service
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the service
//     schema:
//       "$ref": "#/definitions/Error"

// GetService represents the API handler to capture a
// service for a build from the configured backend.
func GetService(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"org":     o,
		"repo":    r.GetName(),
		"service": s.GetNumber(),
		"user":    u.GetName(),
	}).Infof("reading service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	c.JSON(http.StatusOK, s)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/services/{service} services UpdateService
//
// Update a service for a build in the configured backend
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
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: service
//   description: Service number
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing the service to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Service"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the service
//     schema:
//       "$ref": "#/definitions/Service"
//   '400':
//     description: Unable to update the service
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the service
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateService represents the API handler to update
// a service for a build in the configured backend.
func UpdateService(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"org":     o,
		"repo":    r.GetName(),
		"service": s.GetNumber(),
		"user":    u.GetName(),
	}).Infof("updating service %s", entry)

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update service fields if provided
	if len(input.GetStatus()) > 0 {
		// update status if set
		s.SetStatus(input.GetStatus())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		s.SetError(input.GetError())
	}

	if input.GetExitCode() > 0 {
		// update exit_code if set
		s.SetExitCode(input.GetExitCode())
	}

	if input.GetStarted() > 0 {
		// update started if set
		s.SetStarted(input.GetStarted())
	}

	if input.GetFinished() > 0 {
		// update finished if set
		s.SetFinished(input.GetFinished())
	}

	// send API call to update the service
	err = database.FromContext(c).UpdateService(s)
	if err != nil {
		retErr := fmt.Errorf("unable to update service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated service
	s, _ = database.FromContext(c).GetService(s.GetNumber(), b)

	c.JSON(http.StatusOK, s)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/services/{service} services DeleteService
//
// Delete a service for a build in the configured backend
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
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: service
//   description: Service Number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the service
//     schema:
//       type: string
//   '500':
//     description: Unable to delete the service
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteService represents the API handler to remove
// a service for a build from the configured backend.
//
// nolint: dupl // ignore similar code with step
func DeleteService(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"org":     o,
		"repo":    r.GetName(),
		"service": s.GetNumber(),
		"user":    u.GetName(),
	}).Infof("deleting service %s", entry)

	// send API call to remove the service
	err := database.FromContext(c).DeleteService(s.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("service %s deleted", entry))
}

// planServices is a helper function to plan all services
// in the build for execution. This creates the services
// for the build in the configured backend.
//
// nolint: lll // ignore long line length due to variable names
func planServices(database database.Service, p *pipeline.Build, b *library.Build) ([]*library.Service, error) {
	// variable to store planned services
	services := []*library.Service{}

	// iterate through all pipeline services
	for _, service := range p.Services {
		// create the service object
		s := new(library.Service)
		s.SetBuildID(b.GetID())
		s.SetRepoID(b.GetRepoID())
		s.SetName(service.Name)
		s.SetImage(service.Image)
		s.SetNumber(service.Number)
		s.SetStatus(constants.StatusPending)
		s.SetCreated(time.Now().UTC().Unix())

		// send API call to create the service
		err := database.CreateService(s)
		if err != nil {
			return services, fmt.Errorf("unable to create service %s: %w", s.GetName(), err)
		}

		// send API call to capture the created service
		s, err = database.GetService(s.GetNumber(), b)
		if err != nil {
			return services, fmt.Errorf("unable to get service %s: %w", s.GetName(), err)
		}

		// populate environment variables from service library
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#Service.Environment
		err = service.MergeEnv(s.Environment())
		if err != nil {
			return services, err
		}

		// create the log object
		l := new(library.Log)
		l.SetServiceID(s.GetID())
		l.SetBuildID(b.GetID())
		l.SetRepoID(b.GetRepoID())
		l.SetData([]byte{})

		// send API call to create the service logs
		err = database.CreateLog(l)
		if err != nil {
			return services, fmt.Errorf("unable to create service logs for service %s: %w", s.GetName(), err)
		}
	}

	return services, nil
}
