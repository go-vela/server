// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/steps steps CreateStep
//
// Create a step for a build
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
//   description: Payload containing the step to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Step"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the step
//     schema:
//       "$ref": "#/definitions/Step"
//   '400':
//     description: Unable to create the step
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the step
//     schema:
//       "$ref": "#/definitions/Error"

// CreateStep represents the API handler to create
// a step for a build in the configured backend.
func CreateStep(c *gin.Context) {
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
	}).Infof("creating new step for build %s", entry)

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new step for build %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in step object
	input.SetRepoID(r.GetID())
	input.SetBuildID(b.GetID())

	if len(input.GetStatus()) == 0 {
		input.SetStatus(constants.StatusPending)
	}

	if input.GetCreated() == 0 {
		input.SetCreated(time.Now().UTC().Unix())
	}

	// send API call to create the step
	err = database.FromContext(c).CreateStep(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create step for build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created step
	s, _ := database.FromContext(c).GetStepForBuild(b, input.GetNumber())

	c.JSON(http.StatusCreated, s)
}

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/steps steps GetSteps
//
// Retrieve a list of steps for a build
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
//     description: Successfully retrieved the list of steps
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Step"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of steps
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of steps
//     schema:
//       "$ref": "#/definitions/Error"

// GetSteps represents the API handler to capture a list
// of steps for a build from the configured backend.
func GetSteps(c *gin.Context) {
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
	}).Infof("reading steps for build %s", entry)

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
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of steps for the build
	s, t, err := database.FromContext(c).ListStepsForBuild(b, map[string]interface{}{}, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get steps for build %s: %w", entry, err)

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

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step} steps GetStep
//
// Retrieve a step for a build
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
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the step
//     schema:
//       "$ref": "#/definitions/Step"

// GetStep represents the API handler to capture a
// step for a build from the configured backend.
func GetStep(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"step":  s.GetNumber(),
		"user":  u.GetName(),
	}).Infof("reading step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	c.JSON(http.StatusOK, s)
}

// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step} steps UpdateStep
//
// Update a step for a build
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
// - in: body
//   name: body
//   description: Payload containing the step to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Step"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the step
//     schema:
//       "$ref": "#/definitions/Step"
//   '400':
//     description: Unable to update the step
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the step
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateStep represents the API handler to update
// a step for a build in the configured backend.
func UpdateStep(c *gin.Context) {
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
	}).Infof("updating step %s", entry)

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update step fields if provided
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

	if len(input.GetHost()) > 0 {
		// update host if set
		s.SetHost(input.GetHost())
	}

	if len(input.GetRuntime()) > 0 {
		// update runtime if set
		s.SetRuntime(input.GetRuntime())
	}

	if len(input.GetDistribution()) > 0 {
		// update distribution if set
		s.SetDistribution(input.GetDistribution())
	}

	// send API call to update the step
	err = database.FromContext(c).UpdateStep(s)
	if err != nil {
		retErr := fmt.Errorf("unable to update step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated step
	s, _ = database.FromContext(c).GetStepForBuild(b, s.GetNumber())

	c.JSON(http.StatusOK, s)
}

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/steps/{step} steps DeleteStep
//
// Delete a step for a build
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
//     description: Successfully deleted the step
//     schema:
//       type: string
//   '500':
//     description: Successfully deleted the step
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteStep represents the API handler to remove
// a step for a build from the configured backend.
func DeleteStep(c *gin.Context) {
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
	}).Infof("deleting step %s", entry)

	// send API call to remove the step
	err := database.FromContext(c).DeleteStep(s)
	if err != nil {
		retErr := fmt.Errorf("unable to delete step %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("step %s deleted", entry))
}

// planSteps is a helper function to plan all steps
// in the build for execution. This creates the steps
// for the build in the configured backend.
func planSteps(database database.Interface, p *pipeline.Build, b *library.Build) ([]*library.Step, error) {
	// variable to store planned steps
	steps := []*library.Step{}

	// iterate through all pipeline stages
	for _, stage := range p.Stages {
		// iterate through all steps for each pipeline stage
		for _, step := range stage.Steps {
			// create the step object
			s := new(library.Step)
			s.SetBuildID(b.GetID())
			s.SetRepoID(b.GetRepoID())
			s.SetNumber(step.Number)
			s.SetName(step.Name)
			s.SetImage(step.Image)
			s.SetStage(stage.Name)
			s.SetStatus(constants.StatusPending)
			s.SetCreated(time.Now().UTC().Unix())

			// send API call to create the step
			err := database.CreateStep(s)
			if err != nil {
				return steps, fmt.Errorf("unable to create step %s: %w", s.GetName(), err)
			}

			// send API call to capture the created step
			s, err = database.GetStepForBuild(b, s.GetNumber())
			if err != nil {
				return steps, fmt.Errorf("unable to get step %s: %w", s.GetName(), err)
			}

			// populate environment variables from step library
			//
			// https://pkg.go.dev/github.com/go-vela/types/library#step.Environment
			err = step.MergeEnv(s.Environment())
			if err != nil {
				return steps, err
			}

			// create the log object
			l := new(library.Log)
			l.SetStepID(s.GetID())
			l.SetBuildID(b.GetID())
			l.SetRepoID(b.GetRepoID())
			l.SetData([]byte{})

			// send API call to create the step logs
			err = database.CreateLog(l)
			if err != nil {
				return nil, fmt.Errorf("unable to create logs for step %s: %w", s.GetName(), err)
			}

			steps = append(steps, s)
		}
	}

	// iterate through all pipeline steps
	for _, step := range p.Steps {
		// create the step object
		s := new(library.Step)
		s.SetBuildID(b.GetID())
		s.SetRepoID(b.GetRepoID())
		s.SetNumber(step.Number)
		s.SetName(step.Name)
		s.SetImage(step.Image)
		s.SetStatus(constants.StatusPending)
		s.SetCreated(time.Now().UTC().Unix())

		// send API call to create the step
		err := database.CreateStep(s)
		if err != nil {
			return steps, fmt.Errorf("unable to create step %s: %w", s.GetName(), err)
		}

		// send API call to capture the created step
		s, err = database.GetStepForBuild(b, s.GetNumber())
		if err != nil {
			return steps, fmt.Errorf("unable to get step %s: %w", s.GetName(), err)
		}

		// populate environment variables from step library
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#step.Environment
		err = step.MergeEnv(s.Environment())
		if err != nil {
			return steps, err
		}

		// create the log object
		l := new(library.Log)
		l.SetStepID(s.GetID())
		l.SetBuildID(b.GetID())
		l.SetRepoID(b.GetRepoID())
		l.SetData([]byte{})

		// send API call to create the step logs
		err = database.CreateLog(l)
		if err != nil {
			return steps, fmt.Errorf("unable to create logs for step %s: %w", s.GetName(), err)
		}

		steps = append(steps, s)
	}

	return steps, nil
}
