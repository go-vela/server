// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

const (
	// ScheduleResp represents a JSON return for a single schedule.
	ScheduleResp = `{
		"id": 2,
		"repo": {
			"id": 1,
			"owner": {
				"id": 1,
				"name": "octocat",
				"favorites": [],
				"active": true,
				"admin": false
			},
			"org": "github",
			"counter": 10,
			"name": "octocat",
			"full_name": "github/octocat",
			"link": "https://github.com/github/octocat",
			"clone": "https://github.com/github/octocat",
			"branch": "main",
			"build_limit": 10,
			"timeout": 60,
			"visibility": "public",
			"private": false,
			"trusted": true,
			"pipeline_type": "yaml",
			"topics": [],
			"active": true,
			"allow_events": {
				"push": {
					"branch": true,
					"tag": true
				},
				"pull_request": {
					"opened": true,
					"synchronize": true,
					"reopened": true,
					"edited": false
				},
				"deployment": {
					"created": true
				},
				"comment": {
					"created": false,
					"edited": false
				}
			},
			"approve_build": "fork-always",
			"previous_name": ""
		},
		"active": true,
		"name": "foo",
		"entry": "@weekly",
		"created_at": 1683154980,
		"created_by": "octocat",
		"updated_at": 1683154980,
		"updated_by": "octocat",
		"scheduled_at": 0,
		"branch": "main",
		"error": "error message",
		"next_run": 0
	}`
	SchedulesResp = `[
	{
		"id": 2,
		"repo": {
			"id": 1,
			"owner": {
				"id": 1,
				"name": "octocat",
				"favorites": [],
				"active": true,
				"admin": false
			},
			"org": "github",
			"counter": 10,
			"name": "octocat",
			"full_name": "github/octocat",
			"link": "https://github.com/github/octocat",
			"clone": "https://github.com/github/octocat",
			"branch": "main",
			"build_limit": 10,
			"timeout": 60,
			"visibility": "public",
			"private": false,
			"trusted": true,
			"pipeline_type": "yaml",
			"topics": [],
			"active": true,
			"allow_events": {
				"push": {
					"branch": true,
					"tag": true
				},
				"pull_request": {
					"opened": true,
					"synchronize": true,
					"reopened": true,
					"edited": false
				},
				"deployment": {
					"created": true
				},
				"comment": {
					"created": false,
					"edited": false
				}
			},
			"approve_build": "fork-always",
			"previous_name": ""
		},
		"active": true,
		"name": "foo",
		"entry": "@weekly",
		"created_at": 1683154980,
		"created_by": "octocat",
		"updated_at": 1683154980,
		"updated_by": "octocat",
		"scheduled_at": 0,
		"branch": "main",
		"error": "error message",
		"next_run": 0
	},
	{
		"id": 1,
		"repo": {
			"id": 1,
			"owner": {
				"id": 1,
				"name": "octocat",
				"favorites": [],
				"active": true,
				"admin": false
			},
			"org": "github",
			"counter": 10,
			"name": "octocat",
			"full_name": "github/octocat",
			"link": "https://github.com/github/octocat",
			"clone": "https://github.com/github/octocat",
			"branch": "main",
			"build_limit": 10,
			"timeout": 60,
			"visibility": "public",
			"private": false,
			"trusted": true,
			"pipeline_type": "yaml",
			"topics": [],
			"active": true,
			"allow_events": {
				"push": {
					"branch": true,
					"tag": true
				},
				"pull_request": {
					"opened": true,
					"synchronize": true,
					"reopened": true,
					"edited": false
				},
				"deployment": {
					"created": true
				},
				"comment": {
					"created": false,
					"edited": false
				}
			},
			"approve_build": "fork-always",
			"previous_name": ""
		},
		"active": true,
		"name": "bar",
		"entry": "@weekly",
		"created_at": 1683154974,
		"created_by": "octocat",
		"updated_at": 1683154974,
		"updated_by": "octocat",
		"scheduled_at": 0,
		"repo_id": 1,
		"branch": "main",
		"error": "error message",
		"next_run": 0
	}]`
)

// getSchedules returns mock JSON for a http GET.
func getSchedules(c *gin.Context) {
	data := []byte(SchedulesResp)

	var body []api.Schedule

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getSchedule has a param :schedule returns mock JSON for a http GET.
//
// Pass "not-found" to :schedule to test receiving a http 404 response.
func getSchedule(c *gin.Context) {
	s := c.Param("schedule")

	if strings.Contains(s, "not-found") {
		msg := fmt.Sprintf("Schedule %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	data := []byte(ScheduleResp)

	var body api.Schedule

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addSchedule returns mock JSON for a http POST.
func addSchedule(c *gin.Context) {
	data := []byte(ScheduleResp)

	var body api.Schedule

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateSchedule has a param :schedule returns mock JSON for a http PUT.
//
// Pass "not-found" to :schedule to test receiving a http 404 response.
func updateSchedule(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		s := c.Param("schedule")

		if strings.Contains(s, "not-found") {
			msg := fmt.Sprintf("Schedule %s does not exist", s)

			c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

			return
		}
	}

	data := []byte(ScheduleResp)

	var body api.Schedule

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeSchedule has a param :schedule returns mock JSON for a http DELETE.
//
// Pass "not-found" to :schedule to test receiving a http 404 response.
func removeSchedule(c *gin.Context) {
	s := c.Param("schedule")

	if strings.Contains(s, "not-found") {
		msg := fmt.Sprintf("Schedule %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("schedule %s deleted", s))
}
