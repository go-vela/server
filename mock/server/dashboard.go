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
	// DashboardResp represents a JSON return for a dashboard.
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

	// DashCardResp represents a JSON return for a DashCard.
	DashCardResp = `{
    "dashboard": {
        "id": "6e9f84c3-d853-4afb-b56e-99ff200264c0",
        "name": "dashboard-1",
        "created_at": 1714677999,
        "created_by": "Octocat",
        "updated_at": 1714678173,
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
                "id": 2,
                "name": "Octocat/test-repo"
            },
            {
                "id": 1,
                "name": "Octocat/test-repo-2"
            }
        ]
    },
    "repos": [
        {
            "org": "Octocat",
            "name": "test-repo",
            "counter": 1,
            "builds": [
                {
                    "number": 1,
                    "started": 1714678666,
                    "finished": 1714678672,
                    "sender": "Octocat",
                    "status": "failure",
                    "event": "deployment",
                    "branch": "refs/heads/main",
                    "link": "http://vela/Octocat/test-repo/1"
                }
            ]
        },
        {
            "org": "Octocat",
            "name": "test-repo-2"
        }
    ]
}`

	// DashCardsResp represents a JSON return for multiple DashCards.
	DashCardsResp = `[
{
    "dashboard": {
        "id": "6e9f84c3-d853-4afb-b56e-99ff200264c0",
        "name": "dashboard-1",
        "created_at": 1714677999,
        "created_by": "Octocat",
        "updated_at": 1714678173,
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
                "id": 2,
                "name": "Octocat/test-repo"
            },
            {
                "id": 1,
                "name": "Octocat/test-repo-2"
            }
        ]
    },
    "repos": [
        {
            "org": "Octocat",
            "name": "test-repo",
            "counter": 1,
            "builds": [
                {
                    "number": 1,
                    "started": 1714678666,
                    "finished": 1714678672,
                    "sender": "Octocat",
                    "status": "failure",
                    "event": "deployment",
                    "branch": "refs/heads/main",
                    "link": "http://vela/Octocat/test-repo/1"
                }
            ]
        },
        {
            "org": "Octocat",
            "name": "test-repo-2"
        }
    ]
},
{
    "dashboard": {
        "id": "6e9f84c3-d853-4afb-b56e-99ff200264c1",
        "name": "dashboard-2",
        "created_at": 1714677999,
        "created_by": "Octocat",
        "updated_at": 1714678173,
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
                "id": 2,
                "name": "Octocat/test-repo"
            },
            {
                "id": 1,
                "name": "Octocat/test-repo-2"
            }
        ]
    },
    "repos": [
        {
            "org": "Octocat",
            "name": "test-repo",
            "counter": 1,
            "builds": [
                {
                    "number": 1,
                    "started": 1714678666,
                    "finished": 1714678672,
                    "sender": "Octocat",
                    "status": "failure",
                    "event": "deployment",
                    "branch": "refs/heads/main",
                    "link": "http://vela/Octocat/test-repo/1"
                }
            ]
        },
        {
            "org": "Octocat",
            "name": "test-repo-2"
        }
    ]
}
]`
)

// getDashboards returns mock JSON for a http GET.
func getDashboards(c *gin.Context) {
	data := []byte(DashCardsResp)

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

	data := []byte(DashCardResp)

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
