// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apiTypes "github.com/go-vela/server/api/types"
	"github.com/go-vela/types"
)

const (
	// ScheduleResp represents a JSON return for a single schedule.
	ScheduleResp = `{
  "id": 1,
  "user_id": 1,
  "org": "github",
  "name": "octocat",
  "full_name": "github/octocat",
  "link": "https://github.com/github/octocat",
  "clone": "https://github.com/github/octocat",
  "branch": "master",
  "build_limit": 10,
  "timeout": 60,
  "visibility": "public",
  "private": false,
  "trusted": true,
  "active": true,
  "allow_pr": false,
  "allow_push": true,
  "allow_deploy": false,
  "allow_tag": false
}`

	// SchedulesResp represents a JSON return for one to many schedules.
	SchedulesResp = `[
  {
    "id": 1,
    "user_id": 1,
    "org": "github",
    "name": "octocat",
    "full_name": "github/octocat",
    "link": "https://github.com/github/octocat",
    "clone": "https://github.com/github/octocat",
    "branch": "master",
    "build_limit": 10,
    "timeout": 60,
    "visibility": "public",
    "private": false,
    "trusted": true,
    "active": true,
    "allow_pr": false,
    "allow_push": true,
    "allow_deploy": false,
    "allow_tag": false
  },
  {
    "id": 2,
    "user_id": 1,
    "org": "github",
    "name": "octokitty",
    "full_name": "github/octokitty",
    "link": "https://github.com/github/octokitty",
    "clone": "https://github.com/github/octokitty",
    "branch": "master",
    "build_limit": 10,
    "timeout": 60,
    "visibility": "public",
    "private": false,
    "trusted": true,
    "active": true,
    "allow_pr": false,
    "allow_push": true,
    "allow_deploy": false,
    "allow_tag": false
  }
]`
)

// getSchedules returns mock JSON for a http GET.
func getSchedules(c *gin.Context) {
	data := []byte(SchedulesResp)

	var body []apiTypes.Schedule
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

	var body apiTypes.Schedule
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addSchedule returns mock JSON for a http POST.
func addSchedule(c *gin.Context) {
	data := []byte(ScheduleResp)

	var body apiTypes.Schedule
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

	var body apiTypes.Schedule
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
