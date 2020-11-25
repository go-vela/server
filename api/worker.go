// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/workers workers CreateWorker
//
// Create a worker for the configured backend
//
// ---
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the worker to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Worker"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the worker
//     type: json
//     schema:
//       "$ref": "#/definitions/Worker"
//   '400':
//     description: Unable to create the worker
//     schema:
//       type: string
//   '500':
//     description: Unable to create the worker
//     schema:
//       type: string

// CreateWorker represents the API handler to
// create a worker in the configured backend
func CreateWorker(c *gin.Context) {
	logrus.Info("Creating new worker")

	// capture body from API request
	input := new(library.Worker)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new worker: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	err = database.FromContext(c).CreateWorker(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create worker: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, fmt.Sprintf("worker %s created", input.GetHostname()))
}

// swagger:operation GET /api/v1/workers workers GetWorkers
//
// Retrieve a list of workers for the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the list of workers
//     type: json
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Worker"
//   '500':
//     description: Unable to retrieve the list of workers
//     schema:
//       type: string

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

// swagger:operation GET /api/v1/workers/{worker} workers GetWorker
//
// Retrieve a worker for the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: Hostname of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the worker
//     type: json
//     schema:
//       "$ref": "#/definitions/Worker"
//   '404':
//     description: Unable to retrieve the worker
//     schema:
//       type: string

// GetWorker represents the API handler to capture a
// worker from the configured backend.
func GetWorker(c *gin.Context) {
	w := worker.Retrieve(c)
	w, err := database.FromContext(c).GetWorker(w.GetHostname())
	if err != nil {
		retErr := fmt.Errorf("unable to get workers: %w", err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	c.JSON(http.StatusOK, w)
}

// swagger:operation PUT /api/v1/workers/{worker} workers UpdateWorker
//
// Update a worker for the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the worker to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Worker"
// - in: path
//   name: worker
//   description: Name of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the worker
//     type: json
//     schema:
//       "$ref": "#/definitions/Worker"
//   '400':
//     description: Unable to update the worker
//     schema:
//       type: string
//   '404':
//     description: Unable to update the worker
//     schema:
//       type: string
//   '500':
//     description: Unable to update the worker
//     schema:
//       type: string

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

	if input.GetLastCheckedIn() > 0 {
		// update LastCheckedIn if set
		w.SetLastCheckedIn(input.GetLastCheckedIn())
	}

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

// swagger:operation DELETE /api/v1/workers/{worker} workers DeleteWorker
//
// Delete a worker for the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: Name of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted of worker
//     schema:
//       type: string
//   '404':
//     description: Unable to delete worker
//     schema:
//       type: string
//   '500':
//     description: Unable to delete worker
//     schema:
//       type: string

// DeleteWorker represents the API handler to remove
// a worker from the configured backend.
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
