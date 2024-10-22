// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/router/middleware/auth"
)

const (
	// UserResp represents a JSON return for a single user.
	UserResp = `{
  "id": 1,
  "name": "OctoKitty",
  "token": null,
  "favorites": ["github/octocat"],
  "active": true,
  "admin": false,
  "dashboards": []
}`

	// UsersResp represents a JSON return for one to many users.
	UsersResp = `[
  {
    "id": 2,
    "name": "octocat",
    "token": null,
    "favorites": ["github/octocat"],
    "active": true,
    "admin": false,
	"dashboards": []
  },
  {
    "id": 1,
    "name": "OctoKitty",
    "token": null,
    "favorites": ["github/octocat"],
    "active": true,
    "admin": false,
	"dashboards": []
  }
]`
)

// getUsers returns mock JSON for a http GET.
func getUsers(c *gin.Context) {
	data := []byte(UsersResp)

	var body []api.User
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getUser has a param :user returns mock JSON for a http GET.
//
// Pass "not-found" to :user to test receiving a http 404 response.
func getUser(c *gin.Context) {
	u := c.Param("user")

	if strings.Contains(u, "not-found") {
		msg := fmt.Sprintf("User %s does not exist", u)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	data := []byte(UserResp)

	var body api.User
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// currentUser returns mock JSON for a http GET and http PUT.
//
// Pass "invalid" to auth header to test receiving 401 response.
func currentUser(c *gin.Context) {
	tkn, _ := auth.RetrieveAccessToken(c.Request)

	if strings.Contains(tkn, "invalid") {
		msg := "unauthorized"

		c.AbortWithStatusJSON(http.StatusUnauthorized, api.Error{Message: &msg})

		return
	}

	data := []byte(UserResp)

	var body api.User
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addUser returns mock JSON for a http POST.
func addUser(c *gin.Context) {
	data := []byte(UserResp)

	var body api.User
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateUser has a param :user returns mock JSON for a http PUT.
//
// Pass "not-found" to :user to test receiving a http 404 response.
func updateUser(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		u := c.Param("user")

		if strings.Contains(u, "not-found") {
			msg := fmt.Sprintf("User %s does not exist", u)

			c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

			return
		}
	}

	data := []byte(UserResp)

	var body api.User
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeUser has a param :user returns mock JSON for a http DELETE.
//
// Pass "not-found" to :user to test receiving a http 404 response.
func removeUser(c *gin.Context) {
	u := c.Param("user")

	if strings.Contains(u, "not-found") {
		msg := fmt.Sprintf("User %s does not exist", u)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("User %s removed", u))
}
