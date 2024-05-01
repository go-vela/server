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
)

const (
	// DashboardResp represents a JSON return for a single build.
	DashboardResp = `{
  "id": "c976470d-34c1-49b2-9a98-1035871c576b",
  "name": "my-dashboard",
  "created_at": 1714573212,
  "created_by": "Octocat",
  "updated_at": 1714573212,
  "updated_by": "Octocat",
  "admins": [
    {
      "id": 1,
      "name": "Octocat",
      "active": true
    }
  ],
  "repos": [
    {
      "id": 1,
      "name": "Octocat/vela-repo",
      "branches": [
        "main"
      ],
      "events": [
        "push"
      ]
    }
  ]
}`

	// DashboardsResp represents a JSON return for one to many builds.
	DashboardsResp = `[
{
  "id": "c976470d-34c1-49b2-9a98-1035871c576b",
  "name": "my-dashboard",
  "created_at": 1714573212,
  "created_by": "Octocat",
  "updated_at": 1714573212,
  "updated_by": "Octocat",
  "admins": [
    {
      "id": 1,
      "name": "Octocat",
      "active": true
    }
  ],
  "repos": [
    {
      "id": 1,
      "name": "Octocat/vela-repo",
      "branches": [
        "main"
      ],
      "events": [
        "push"
      ]
    }
  ]
},
{
  "id": "c976470d-34c1-49b2-9a98-1035871c576c",
  "name": "my-second-dashboard",
  "created_at": 1714573212,
  "created_by": "Octocat",
  "updated_at": 1714573212,
  "updated_by": "Octocat",
  "admins": [
    {
      "id": 1,
      "name": "Octocat",
      "active": true
    }
  ],
  "repos": [
    {
      "id": 1,
      "name": "Octocat/vela-repo",
      "branches": [
        "main"
      ],
      "events": [
        "push"
      ]
    }
  ]
}
]`
)

// getDashboards returns mock JSON for a http GET.
func getDashboards(c *gin.Context) {
	data := []byte(DashboardsResp)

	var body []api.Dashboard
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getDashboard has a param :dashboard returns mock JSON for a http GET.
func getDashboard(c *gin.Context) {
	d := c.Param("dashboard")

	if strings.EqualFold(d, "0") {
		msg := fmt.Sprintf("Dashboard %s does not exist", d)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(DashboardResp)

	var body api.Dashboard
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addDashboard returns mock JSON for a http POST.
func addDashboard(c *gin.Context) {
	data := []byte(DashboardResp)

	var body api.Dashboard
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateDashboard returns mock JSON for a http PUT.
func updateDashboard(c *gin.Context) {
	data := []byte(DashboardResp)

	var body api.Dashboard
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
