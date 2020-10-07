// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateWorker represents the API handler to
// create a worker in the configured backend
func CreateWorker(c *gin.Context) {
	input := new(library.Worker)

	err := c.Bind(input)

	// set LastCheckedIn to now
	input.SetLastCheckedIn(time.Now().Unix())

	err = database.FromContext(c).CreateWorker(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create worker: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("worker %s created", input.GetHostname()))
}

// UpdateWorker represents the API handler to
// create a worker in the configured backend
func UpdateWorker(c *gin.Context) {
	// capture middleware values
	worker := c.Param("worker")

	logrus.Infof("Updating worker %s", worker)

	// capture body from API request
	input := new(library.Worker)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for worker %s: %w", worker, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the worker
	w, err := database.FromContext(c).GetWorker(worker)
	if err != nil {
		retErr := fmt.Errorf("unable to get worker %s: %w", worker, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	if len(input.GetAddress()) > 0 {
		// update admin if set
		w.SetAddress(input.GetAddress())
	}

	if len(input.GetRoutes()) > 0 {
		// update routes if set
		w.SetRoutes(input.GetRoutes())
	}

	if input.GetActive() {
		// update active if set
		w.SetActive(input.GetActive())
	}

	// update LastCheckedIn to now
	w.SetLastCheckedIn(time.Now().Unix())

	// send API call to update the worker
	err = database.FromContext(c).UpdateWorker(w)
	if err != nil {
		retErr := fmt.Errorf("unable to update worker %s: %w", worker, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated worker
	w, _ = database.FromContext(c).GetWorker(worker)

	c.JSON(http.StatusOK, w)
}

// GetWorkers represents the API handler to capture a
// list of workers from the configured backend.
func GetWorkers(c *gin.Context) {
	w, err := database.FromContext(c).GetWorkerList()
	if err != nil {
		retErr := fmt.Errorf("unable to get workers: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, w)
}

// GetWorker represents the API handler to capture a
// list of workers from the configured backend.
func GetWorker(c *gin.Context) {
	w := worker.Retrieve(c)
	w, err := database.FromContext(c).GetWorker(w.GetHostname())
	if err != nil {
		retErr := fmt.Errorf("unable to get workers: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, w)
}

// DeleteWorker represents the API handler to remove
// a worker for a build from the configured backend.
func DeleteWorker(c *gin.Context) {
	// capture middleware values
	w := worker.Retrieve(c)

	logrus.Infof("Deleting worker %s", w.GetHostname())

	// send API call to remove the step
	err := database.FromContext(c).DeleteWorker(w.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Worker %s deleted", w.GetHostname()))
}
