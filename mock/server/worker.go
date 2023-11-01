// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
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
			"running_builds": [
  				{
					"id": 2,
					"repo_id": 1,
					"number": 2,
					"parent": 1,
					"event": "push",
					"status": "running",
					"error": "",
					"enqueued": 1563474204,
					"created": 1563474204,
					"started": 1563474204,
					"finished": 0,
					"deploy": "",
					"clone": "https://github.com/github/octocat.git",
					"source": "https://github.com/github/octocat/commit/48afb5bdc41ad69bf22588491333f7cf71135163",
					"title": "push received from https://github.com/github/octocat",
					"message": "Second commit...",
					"commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
					"sender": "OctoKitty",
					"author": "OctoKitty",
					"email": "octokitty@github.com",
					"link": "https://vela.example.company.com/github/octocat/1",
					"branch": "main",
					"ref": "refs/heads/main",
					"base_ref": "",
					"host": "ed95dcc0687c",
					"runtime": "",
					"distribution": ""
  				},
  				{
					"id": 1,
					"repo_id": 1,
					"number": 1,
					"parent": 1,
					"event": "push",
					"status": "running",
					"error": "",
					"enqueued": 1563474077,
					"created": 1563474076,
					"started": 1563474077,
					"finished": 0,
					"deploy": "",
					"clone": "https://github.com/github/octocat.git",
					"source": "https://github.com/github/octocat/commit/48afb5bdc41ad69bf22588491333f7cf71135163",
					"title": "push received from https://github.com/github/octocat",
					"message": "First commit...",
					"commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
					"sender": "OctoKitty",
					"author": "OctoKitty",
					"email": "octokitty@github.com",
					"link": "https://vela.example.company.com/github/octocat/1",
					"branch": "main",
					"ref": "refs/heads/main",
					"base_ref": "",
					"host": "82823eb770b0",
					"runtime": "",
					"distribution": ""
  				}
			]
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

	// AddWorkerResp represents a JSON return for adding a worker.
	AddWorkerResp = `{
    	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ3b3JrZXIiLCJpYXQiOjE1MTYyMzkwMjIsInRva2VuX3R5cGUiOiJXb3JrZXJBdXRoIn0.qeULIimCJlrwsE0JykNpzBmMaHUbvfk0vkyAz2eEo38"
    }`

	// RefreshWorkerAuthResp represents a JSON return for refreshing a worker's authentication.
	RefreshWorkerAuthResp = `{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ3b3JrZXIiLCJpYXQiOjE1MTYyMzkwMjIsInRva2VuX3R5cGUiOiJXb3JrZXJBdXRoIn0.qeULIimCJlrwsE0JykNpzBmMaHUbvfk0vkyAz2eEo38"
	}`

	// RegisterTokenResp represents a JSON return for an admin requesting a registration token.
	//
	//nolint:gosec // not actual credentials
	RegisterTokenResp = `{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ3b3JrZXIiLCJpYXQiOjE1MTYyMzkwMjIsInRva2VuX3R5cGUiOiJXb3JrZXJSZWdpc3RlciJ9.gEzKaZB-sDd_gFCVF5uGo2mcf3iy9CrXDTLPZ6PTsTc"
	}`

	// QueueInfoResp represents a JSON return for an admin requesting a queue registration info.
	//
	//not actual credentials.
	QueueInfoResp = `{
		"queue_public_key": "DXeyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ98zmko=",
		"queue_address": "redis://redis:6000"
	}`
)

// getWorkers returns mock JSON for a http GET.
func getWorkers(c *gin.Context) {
	data := []byte(WorkersResp)

	var body []api.Worker
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

	var body api.Worker
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addWorker returns mock JSON for a http POST.
func addWorker(c *gin.Context) {
	data := []byte(AddWorkerResp)

	var body library.Token
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

	var body api.Worker
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// refreshWorkerAuth has a param :worker returns mock JSON for a http PUT.
//
// Pass "0" to :worker to test receiving a http 404 response.
func refreshWorkerAuth(c *gin.Context) {
	w := c.Param("worker")

	if strings.EqualFold(w, "0") {
		msg := fmt.Sprintf("Worker %s does not exist", w)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(RefreshWorkerAuthResp)

	var body library.Token
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

// registerToken has a param :worker returns mock JSON for a http POST.
//
// Pass "0" to :worker to test receiving a http 401 response.
func registerToken(c *gin.Context) {
	w := c.Param("worker")

	if strings.EqualFold(w, "0") {
		msg := fmt.Sprintf("user %s is not a platform admin", w)

		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &msg})

		return
	}

	data := []byte(RegisterTokenResp)

	var body library.Token
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// getQueueCreds returns mock JSON for a http GET.
//
// Pass "" to Authorization header to test receiving a http 401 response.
func getQueueCreds(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	// verify token if empty
	if token == "" {
		msg := "unable get queue credentials; invalid registration token"

		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &msg})

		return
	}

	data := []byte(QueueInfoResp)

	var body library.QueueInfo
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}
