// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/logs builds GetBuildLogs
//
// Get logs for a build in the configured backend
//
// ---
// x-success_http_code: '200'
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
//   description: Build number to restart
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved logs for the build
//     type: json
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Log"
//   '500':
//     description: Unable to retrieve logs for the build
//     schema:
//       type: string

// GetBuildLogs represents the API handler to capture a
// list of logs for a build from the configured backend.
func GetBuildLogs(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)

	logrus.Infof("Reading logs for build %s/%d", r.GetFullName(), b.GetNumber())

	// send API call to capture the list of logs for the build
	l, err := database.FromContext(c).GetBuildLogs(b.GetID())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to get logs for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services CreateServiceLogs
//
// Create the logs for a service
//
// ---
// x-success_http_code: '201'
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
//   description: ID of the service
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing the log to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Log"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the service logs
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to create the service logs
//     schema:
//       type: string
//   '500':
//     description: Unable to create the service logs
//     schema:
//       type: string

// CreateServiceLog represents the API handler to create
// the logs for a service in the configured backend.
//
// nolint: dupl // ignore similar code with step
func CreateServiceLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)

	logrus.Infof("Creating logs for service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// capture body from API request
	input := new(library.Log)

	err := c.Bind(input)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to decode JSON for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in log object
	input.SetServiceID(s.GetID())
	input.SetBuildID(b.GetID())
	input.SetRepoID(r.GetID())

	// send API call to create the logs
	err = database.FromContext(c).CreateLog(input)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to create logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created log
	l, _ := database.FromContext(c).GetServiceLog(s.GetID())

	c.JSON(http.StatusCreated, l)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services GetServiceLogs
//
// Retrieve the logs for a service
//
// ---
// x-success_http_code: '200'
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
//   description: ID of the service
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the service logs
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"
//   '500':
//     description: Unable to retrieve the service logs
//     schema:
//       type: string

// GetServiceLog represents the API handler to capture
// the logs for a service from the configured backend.
func GetServiceLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)

	logrus.Infof("Reading logs for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to capture the service logs
	l, err := database.FromContext(c).GetServiceLog(s.GetID())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to get logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services UpdateServiceLog
//
// Update the logs for a service
//
// ---
// x-success_http_code: '201'
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
// - in: body
//   name: body
//   description: Payload containing the log to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Log"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the service logs
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to updated the service logs
//     schema:
//       type: string
//   '500':
//     description: Unable to updates the service logs
//     schema:
//       type: string

// UpdateServiceLog represents the API handler to update
// the logs for a service in the configured backend.
func UpdateServiceLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)

	logrus.Infof("Updating logs for service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to capture the service logs
	l, err := database.FromContext(c).GetServiceLog(s.GetID())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to get logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// capture body from API request
	input := new(library.Log)

	err = c.Bind(input)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to decode JSON for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update log fields if provided
	if len(input.GetData()) > 0 {
		// update data if set
		l.SetData(input.GetData())
	}

	// send API call to update the log
	err = database.FromContext(c).UpdateLog(l)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to update logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated log
	l, _ = database.FromContext(c).GetServiceLog(s.GetID())

	c.JSON(http.StatusOK, l)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services DeleteServiceLogs
//
// Delete the logs for a service
//
// ---
// x-success_http_code: '201'
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
//   description: ID of the service
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully deleted the service logs
//     schema:
//       type: string
//   '500':
//     description: Unable to delete the service logs
//     schema:
//       type: string

// DeleteServiceLog represents the API handler to remove
// the logs for a service from the configured backend.
//
// nolint: dupl // ignore similar code with step
func DeleteServiceLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)

	logrus.Infof("Deleting logs for service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to remove the log
	err := database.FromContext(c).DeleteLog(s.GetID())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to delete logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// nolint: lll // ignore long line length due to return message
	c.JSON(http.StatusOK, fmt.Sprintf("Logs deleted for service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber()))
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps CreateStepLog
//
// Create the logs for a step
//
// ---
// x-success_http_code: '201'
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
//   name: step
//   description: Build number
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the log to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Log"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the logs for step
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to create the logs for a step
//     schema:
//       type: string
//   '500':
//     description: Unable to create the logs for a step
//     schema:
//       type: string

// CreateStepLog represents the API handler to create
// the logs for a step in the configured backend.
//
// nolint: dupl // ignore similar code with service
func CreateStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	logrus.Infof("Creating logs for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// capture body from API request
	input := new(library.Log)

	err := c.Bind(input)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to decode JSON for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in log object
	input.SetStepID(s.GetID())
	input.SetBuildID(b.GetID())
	input.SetRepoID(r.GetID())

	// send API call to create the logs
	err = database.FromContext(c).CreateLog(input)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to create logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created log
	l, _ := database.FromContext(c).GetStepLog(s.GetID())

	c.JSON(http.StatusCreated, l)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps GetStepLog
//
// Retrieve the logs for a step
//
// ---
// x-success_http_code: '200'
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
//   name: step
//   description: Build number
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the logs for step
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"

// GetStepLog represents the API handler to capture
// the logs for a step from the configured backend.
func GetStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	logrus.Infof("Reading logs for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to capture the step logs
	l, err := database.FromContext(c).GetStepLog(s.GetID())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to get logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps UpdateStepLog
//
// Update the logs for a step
//
// ---
// x-success_http_code: '200'
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
//   name: step
//   description: Build number
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the log to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Log"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the logs for step
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to updated the logs for a step
//     schema:
//       type: string
//   '500':
//     description: Unable to updated the logs for a step
//     schema:
//       type: string

// UpdateStepLog represents the API handler to update
// the logs for a step in the configured backend.
func UpdateStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	logrus.Infof("Updating logs for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to capture the step logs
	l, err := database.FromContext(c).GetStepLog(s.GetID())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to get logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// capture body from API request
	input := new(library.Log)

	err = c.Bind(input)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to decode JSON for step %s/%d/%d: %v", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update log fields if provided
	if len(input.GetData()) > 0 {
		// update data if set
		l.SetData(input.GetData())
	}

	// send API call to update the log
	err = database.FromContext(c).UpdateLog(l)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to update logs for step %s/%d/%d: %v", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated log
	l, _ = database.FromContext(c).GetStepLog(s.GetID())

	c.JSON(http.StatusOK, l)
}

// nolint: lll // ignore long line length due to API path
//
// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps DeleteStepLog
//
// Delete the logs for a step
//
// ---
// x-success_http_code: '200'
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
//   name: step
//   description: Build number
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the logs for the step
//     schema:
//       type: string
//   '500':
//     description: Unable to delete the logs for the step
//     schema:
//       type: string

// DeleteStepLog represents the API handler to remove
// the logs for a step from the configured backend.
//
// nolint: dupl // ignore similar code with service
func DeleteStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	logrus.Infof("Deleting logs for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to remove the log
	err := database.FromContext(c).DeleteLog(s.GetID())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to delete logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// nolint: lll // ignore long line length due to return message
	c.JSON(http.StatusOK, fmt.Sprintf("Logs deleted for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber()))
}
