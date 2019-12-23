// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
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
		retErr := fmt.Errorf("unable to get logs for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

// CreateServiceLog represents the API handler to create
// the logs for a service in the configured backend.
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
		retErr := fmt.Errorf("unable to create logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created log
	l, _ := database.FromContext(c).GetServiceLog(s.GetID())

	c.JSON(http.StatusCreated, l)
}

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
		retErr := fmt.Errorf("unable to get logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

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
		retErr := fmt.Errorf("unable to get logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// capture body from API request
	input := new(library.Log)

	err = c.Bind(input)
	if err != nil {
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
		retErr := fmt.Errorf("unable to update logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated log
	l, _ = database.FromContext(c).GetServiceLog(s.GetID())

	c.JSON(http.StatusOK, l)
}

// DeleteServiceLog represents the API handler to remove
// the logs for a service from the configured backend.
func DeleteServiceLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)

	logrus.Infof("Deleting logs for service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to remove the log
	err := database.FromContext(c).DeleteLog(s.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete logs for service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Logs deleted for service %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber()))
}

// CreateStepLog represents the API handler to create
// the logs for a step in the configured backend.
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
		retErr := fmt.Errorf("unable to create logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created log
	l, _ := database.FromContext(c).GetStepLog(s.GetID())

	c.JSON(http.StatusCreated, l)
}

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
		retErr := fmt.Errorf("unable to get logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}

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
		retErr := fmt.Errorf("unable to get logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// capture body from API request
	input := new(library.Log)
	
	err = c.Bind(input)
	if err != nil {
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
		retErr := fmt.Errorf("unable to update logs for step %s/%d/%d: %v", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated log
	l, _ = database.FromContext(c).GetStepLog(s.GetID())

	c.JSON(http.StatusOK, l)
}

// DeleteStepLog represents the API handler to remove
// the logs for a step from the configured backend.
func DeleteStepLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := step.Retrieve(c)

	logrus.Infof("Deleting logs for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// send API call to remove the log
	err := database.FromContext(c).DeleteLog(s.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete logs for step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), s.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Logs deleted for step %s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber()))
}
