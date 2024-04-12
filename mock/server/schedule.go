// SPDX-License-Identifier: Apache-2.0

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
	// ScheduleResp represents a JSON return for a single schedule.
	ScheduleResp = `{
	"id": 2,
	"repo_id": 1,
	"active": true,
	"name": "foo",
	"entry": "@weekly",
	"created_at": 1683154980,
	"created_by": "octocat",
	"updated_at": 1683154980,
	"updated_by": "octocat",
	"scheduled_at": 0,
	"branch": "main"
}`
	SchedulesResp = `[
	{
		"id": 2,
		"active": true,
		"name": "foo",
		"entry": "@weekly",
		"created_at": 1683154980,
		"created_by": "octocat",
		"updated_at": 1683154980,
		"updated_by": "octocat",
		"scheduled_at": 0,
		"repo_id": 1,
		"branch": "main"
	},
	{
		"id": 1,
		"active": true,
		"name": "bar",
		"entry": "@weekly",
		"created_at": 1683154974,
		"created_by": "octocat",
		"updated_at": 1683154974,
		"updated_by": "octocat",
		"scheduled_at": 0,
		"repo_id": 1,
		"branch": "main
	}
]`
)

// getSchedules returns mock JSON for a http GET.
func getSchedules(c *gin.Context) {
	data := []byte(SchedulesResp)

	var body []library.Schedule
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

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(ScheduleResp)

	var body library.Schedule
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addSchedule returns mock JSON for a http POST.
func addSchedule(c *gin.Context) {
	data := []byte(ScheduleResp)

	var body library.Schedule
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

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(ScheduleResp)

	var body library.Schedule
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

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("schedule %s deleted", s))
}
