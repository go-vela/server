// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
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
//   default: 100
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved logs for the build
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Log"
//   '500':
//     description: Unable to retrieve logs for the build
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildLogs represents the API handler to capture a
// list of logs for a build from the configured backend.
func GetBuildLogs(c *gin.Context) {
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
	}).Infof("reading logs for build %s", entry)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		//nolint:lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to convert page query parameter for build %s: %w", entry, err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for build %s: %w", entry, err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of logs for the build
	l, t, err := database.FromContext(c).ListLogsForBuild(b, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for build %s: %w", entry, err)

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

	c.JSON(http.StatusOK, l)
}

//
// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services CreateServiceLogs
//
// Create the logs for a service
//
// ---
// deprecated: true
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
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to create the service logs
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the service logs
//     schema:
//       "$ref": "#/definitions/Error"

// CreateServiceLog represents the API handler to create
// the logs for a service in the configured backend.
//
//nolint:dupl // ignore similar code with step
func CreateServiceLog(c *gin.Context) {
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
	}).Infof("creating logs for service %s", entry)

	// capture body from API request
	input := new(library.Log)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %s: %w", entry, err)

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
		retErr := fmt.Errorf("unable to create logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created log
	l, _ := database.FromContext(c).GetLogForService(s)

	c.JSON(http.StatusCreated, l)
}

//
// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services GetServiceLogs
//
// Retrieve the logs for a service
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
//   description: ID of the service
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the service logs
//     schema:
//       "$ref": "#/definitions/Log"
//   '500':
//     description: Unable to retrieve the service logs
//     schema:
//       "$ref": "#/definitions/Error"

// GetServiceLog represents the API handler to capture
// the logs for a service from the configured backend.
func GetServiceLog(c *gin.Context) {
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
	}).Infof("reading logs for service %s", entry)

	// send API call to capture the service logs
	l, err := database.FromContext(c).GetLogForService(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

//
// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services UpdateServiceLog
//
// Update the logs for a service
//
// ---
// deprecated: true
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
//   description: Payload containing the log to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Log"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the service logs
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to updated the service logs
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to updates the service logs
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateServiceLog represents the API handler to update
// the logs for a service in the configured backend.
func UpdateServiceLog(c *gin.Context) {
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
	}).Infof("updating logs for service %s", entry)

	// send API call to capture the service logs
	l, err := database.FromContext(c).GetLogForService(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// capture body from API request
	input := new(library.Log)

	err = c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %s: %w", entry, err)

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
		retErr := fmt.Errorf("unable to update logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated log
	l, _ = database.FromContext(c).GetLogForService(s)

	c.JSON(http.StatusOK, l)
}

//
// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services DeleteServiceLogs
//
// Delete the logs for a service
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
//   description: ID of the service
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the service logs
//     schema:
//       type: string
//   '500':
//     description: Unable to delete the service logs
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteServiceLog represents the API handler to remove
// the logs for a service from the configured backend.
func DeleteServiceLog(c *gin.Context) {
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
	}).Infof("deleting logs for service %s", entry)

	// send API call to capture the service logs
	l, err := database.FromContext(c).GetLogForService(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to remove the log
	err = database.FromContext(c).DeleteLog(l)
	if err != nil {
		retErr := fmt.Errorf("unable to delete logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("logs deleted for service %s", entry))
}

//
// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps CreateStepLog
//
// Create the logs for a step
//
// ---
// deprecated: true
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
//   description: Step number
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
//     description: Successfully created the logs for step
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to create the logs for a step
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the logs for a step
//     schema:
//       "$ref": "#/definitions/Error"

// CreateStepLog represents the API handler to create
// the logs for a step in the configured backend.
//
//nolint:dupl // ignore similar code with service
func CreateStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"step":  s.GetNumber(),
		"user":  u.GetName(),
	}).Infof("creating logs for step %s", entry)

	// capture body from API request
	input := new(library.Log)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %s: %w", entry, err)

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
		retErr := fmt.Errorf("unable to create logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created log
	l, _ := database.FromContext(c).GetLogForStep(s)

	c.JSON(http.StatusCreated, l)
}

//
// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps GetStepLog
//
// Retrieve the logs for a step
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
//   name: step
//   description: Step number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the logs for step
//     type: json
//     schema:
//       "$ref": "#/definitions/Log"
//   '500':
//     description: Unable to retrieve the logs for a step
//     schema:
//       "$ref": "#/definitions/Error"

// GetStepLog represents the API handler to capture
// the logs for a step from the configured backend.
func GetStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"step":  s.GetNumber(),
		"user":  u.GetName(),
	}).Infof("reading logs for step %s", entry)

	// send API call to capture the step logs
	l, err := database.FromContext(c).GetLogForStep(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

//
// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps UpdateStepLog
//
// Update the logs for a step
//
// ---
// deprecated: true
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
//   description: Step number
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
//     description: Successfully updated the logs for step
//     schema:
//       "$ref": "#/definitions/Log"
//   '400':
//     description: Unable to update the logs for a step
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the logs for a step
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateStepLog represents the API handler to update
// the logs for a step in the configured backend.
func UpdateStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"step":  s.GetNumber(),
		"user":  u.GetName(),
	}).Infof("updating logs for step %s", entry)

	// send API call to capture the step logs
	l, err := database.FromContext(c).GetLogForStep(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// capture body from API request
	input := new(library.Log)

	err = c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %s: %w", entry, err)

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
		retErr := fmt.Errorf("unable to update logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated log
	l, _ = database.FromContext(c).GetLogForStep(s)

	c.JSON(http.StatusOK, l)
}

//
// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step}/logs steps DeleteStepLog
//
// Delete the logs for a step
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
//   name: step
//   description: Step number
//   required: true
//   type: integer
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
//       "$ref": "#/definitions/Error"

// DeleteStepLog represents the API handler to remove
// the logs for a step from the configured backend.
func DeleteStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"step":  s.GetNumber(),
		"user":  u.GetName(),
	}).Infof("deleting logs for step %s", entry)

	// send API call to capture the step logs
	l, err := database.FromContext(c).GetLogForStep(s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to remove the log
	err = database.FromContext(c).DeleteLog(l)
	if err != nil {
		retErr := fmt.Errorf("unable to delete logs for step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("logs deleted for step %s", entry))
}
