// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreateStep represents the API handler to create
// a step for a build in the configured backend.
func CreateStep(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)

	logrus.Infof("Creating new step for build %s/%d", r.GetFullName(), b.GetNumber())

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new step for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

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
		retErr := fmt.Errorf("unable to create step for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created step
	s, _ := database.FromContext(c).GetStep(input.GetNumber(), b)

	c.JSON(http.StatusCreated, s)
}

// GetSteps represents the API handler to capture a list
// of steps for a build from the configured backend.
func GetSteps(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)

	logrus.Infof("Reading steps for build %s/%d", r.GetFullName(), b.GetNumber())

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the total number of steps for the build
	t, err := database.FromContext(c).GetBuildStepCount(b)
	if err != nil {
		retErr := fmt.Errorf("unable to get steps count for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the list of steps for the build
	s, err := database.FromContext(c).GetBuildStepList(b, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get steps for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

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

// GetStep represents the API handler to capture a
// step for a build from the configured backend.
func GetStep(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)

	logrus.Infof("Reading step %s/%d/%s", r.GetFullName(), b.GetNumber(), c.Param("step"))

	// retrieve step from context
	s := step.Retrieve(c)

	c.JSON(http.StatusOK, s)
}

// UpdateStep represents the API handler to update
// a step for a build in the configured backend.
func UpdateStep(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	logrus.Infof("Updating step %d for build %s/%d", s.GetNumber(), r.GetFullName(), b.GetNumber())

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %s/%d/%d: %v", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

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
		retErr := fmt.Errorf("unable to update step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated step
	s, _ = database.FromContext(c).GetStep(s.GetNumber(), b)

	c.JSON(http.StatusOK, s)
}

// DeleteStep represents the API handler to remove
// a step for a build from the configured backend.
func DeleteStep(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	logrus.Infof("Deleting step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to remove the step
	err := database.FromContext(c).DeleteStep(s.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Step %s/%d/%d deleted", r.GetFullName(), b.GetNumber(), s.GetNumber()))
}

// planSteps is a helper function to plan all steps
// in the build for execution. This creates the steps
// for the build in the configured backend.
func planSteps(database database.Service, p *pipeline.Build, b *library.Build) ([]*library.Step, error) {
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
			s, err = database.GetStep(int(s.GetNumber()), b)
			if err != nil {
				return steps, fmt.Errorf("unable to get step %s: %w", s.GetName(), err)
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
		s, err = database.GetStep(int(s.GetNumber()), b)
		if err != nil {
			return steps, fmt.Errorf("unable to get step %s: %w", s.GetName(), err)
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
