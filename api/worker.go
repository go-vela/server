// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/router/middleware/user"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/workers workers CreateWorker
//
// Create a worker for the configured backend
//
// ---
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
//     schema:
//       type: string
//   '400':
//     description: Unable to create the worker
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the worker
//     schema:
//       "$ref": "#/definitions/Error"

// CreateWorker represents the API handler to
// create a worker in the configured backend.
func CreateWorker(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// capture body from API request
	input := new(library.Worker)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new worker: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user":   u.GetName(),
		"worker": input.GetHostname(),
	}).Infof("creating new worker %s", input.GetHostname())

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
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the list of workers
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Worker"
//   '500':
//     description: Unable to retrieve the list of workers
//     schema:
//       "$ref": "#/definitions/Error"

// GetWorkers represents the API handler to capture a
// list of workers from the configured backend.
func GetWorkers(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Info("reading workers")

	w, err := database.FromContext(c).ListWorkers()
	if err != nil {
		retErr := fmt.Errorf("unable to get workers: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, w)
}

// GetWorkersByStatus represents the API handler to capture a
// list of workers with specified status from the configured backend.
func GetWorkersByStatus(c *gin.Context) {
	s := c.Param("status")
	// capture middleware values
	u := user.Retrieve(c)

	// TODO message/error if not valid status or empty string, or they get back all the workers (GetWorkers), how do other endpoints do it? prob use regex to confirm alpha charas only

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Info("reading workers")

	w, err := database.FromContext(c).ListWorkersByStatus(s)
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
//     schema:
//       "$ref": "#/definitions/Worker"
//   '404':
//     description: Unable to retrieve the worker
//     schema:
//       "$ref": "#/definitions/Error"

// GetWorker represents the API handler to capture a
// worker from the configured backend.
func GetWorker(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	w := worker.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user":   u.GetName(),
		"worker": w.GetHostname(),
	}).Infof("reading worker %s", w.GetHostname())

	w, err := database.FromContext(c).GetWorkerForHostname(w.GetHostname())
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
//     schema:
//       "$ref": "#/definitions/Worker"
//   '400':
//     description: Unable to update the worker
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to update the worker
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the worker
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateWorker represents the API handler to
// create a worker in the configured backend.
func UpdateWorker(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	w := worker.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user":   u.GetName(),
		"worker": w.GetHostname(),
	}).Infof("updating worker %s", w.GetHostname())

	// capture body from API request
	input := new(library.Worker)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

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

	if len(input.GetStatus()) > 0 {
		// update status if set
		w.SetStatus(input.GetStatus())
	}

	if input.GetLastStatusUpdateAt() > 0 {
		// update LastStatusUpdateAt if set
		w.SetLastStatusUpdateAt(input.GetLastStatusUpdateAt())
	}

	if len(input.GetRunningBuildIDs()) > 0 {
		// update RunningBuildIDs if set
		w.SetRunningBuildIDs(input.GetRunningBuildIDs())
	}

	if input.GetLastBuildFinishedAt() > 0 {
		// update LastBuildFinishedAt if set
		w.SetLastBuildFinishedAt(input.GetLastBuildFinishedAt())
	}

	if input.GetLastCheckedIn() > 0 {
		// update LastCheckedIn if set
		w.SetLastCheckedIn(input.GetLastCheckedIn())
	}

	// send API call to update the worker
	err = database.FromContext(c).UpdateWorker(w)
	if err != nil {
		retErr := fmt.Errorf("unable to update worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated worker
	w, _ = database.FromContext(c).GetWorkerForHostname(w.GetHostname())

	c.JSON(http.StatusOK, w)
}

// swagger:operation DELETE /api/v1/workers/{worker} workers DeleteWorker
//
// Delete a worker for the configured backend
//
// ---
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
//   '500':
//     description: Unable to delete worker
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteWorker represents the API handler to remove
// a worker from the configured backend.
func DeleteWorker(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	w := worker.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user":   u.GetName(),
		"worker": w.GetHostname(),
	}).Infof("deleting worker %s", w.GetHostname())

	// send API call to remove the step
	err := database.FromContext(c).DeleteWorker(w)
	if err != nil {
		retErr := fmt.Errorf("unable to delete worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("worker %s deleted", w.GetHostname()))
}
