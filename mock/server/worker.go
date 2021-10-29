// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

const (
	// WorkerResp represents a JSON return for a single worker.
	WorkerResp = `
		{
			"id": 1,
			"hostname": "worker_1",
			"address": "http://vela:8080",
			"routes": [
			"large",
			"docker",
			"large:docker"
			],
			"active": true,
			"last_checked_in": 1602612590
		}`

	// WorkersResp represents a JSON return for one to many workers.
	WorkersResp = `[
		{
			"id": 1,
			"hostname": "worker_1",
			"address": "http://vela:8080",
			"routes": [
			  "large",
			  "docker",
			  "large:docker"
			],
			"active": true,
			"last_checked_in": 1602612590
		  },
		{
			"id": 2,
			"hostname": "worker_2",
			"address": "http://vela:8081",
			"routes": [
			  "large",
			  "docker",
			  "large:docker"
			],
			"active": true,
			"last_checked_in": 1602612590
		  }
	]`
)

// getWorkers returns mock JSON for a http GET.
func getWorkers(c *gin.Context) {
	data := []byte(WorkersResp)

	var body []library.Worker
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getWorker has a param :worker returns mock JSON for a http GET.
func getWorker(c *gin.Context) {
	w := c.Param("worker")

	if strings.EqualFold(w, "0") {
		msg := fmt.Sprintf("Worker %s does not exist", w)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(WorkerResp)

	var body library.Worker
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addWorker returns mock JSON for a http POST.
func addWorker(c *gin.Context) {
	data := []byte(WorkerResp)

	var body library.Worker
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateWorker has a param :worker returns mock JSON for a http PUT.
//
// Pass "0" to :worker to test receiving a http 404 response.
func updateWorker(c *gin.Context) {
	w := c.Param("worker")

	if strings.EqualFold(w, "0") {
		msg := fmt.Sprintf("Worker %s does not exist", w)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(WorkerResp)

	var body library.Worker
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeWorker has a param :worker returns mock JSON for a http DELETE.
//
// Pass "0" to :worker to test receiving a http 404 response.
func removeWorker(c *gin.Context) {
	w := c.Param("worker")

	if strings.EqualFold(w, "0") {
		msg := fmt.Sprintf("Worker %s does not exist", w)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Worker %s removed", w))
}
