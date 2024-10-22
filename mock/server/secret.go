// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore duplicate with user code
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

//nolint:gosec // these are mock responses
const (
	// SecretResp represents a JSON return for a single secret.
	SecretResp = `{
  "id": 1,
  "org": "github",
  "repo": "octocat",
  "team": "",
  "name": "foo",
  "value": "",
  "type": "repo",
  "images": [
    "alpine"
  ],
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
  "allow_command": true,
  "allow_substitution": true,
  "created_at": 1,
  "created_by": "Octocat",
  "updated_at": 2,
  "updated_by": "OctoKitty"
}`

	// SecretsResp represents a JSON return for one to many secrets.
	SecretsResp = `[
  {
    "id": 1,
    "org": "github",
    "repo": "octocat",
    "team": "",
    "name": "foo",
    "value": "",
    "type": "repo",
    "images": [
      "alpine"
    ]
  },
  {
    "id": 2,
    "org": "github",
    "repo": "*",
    "team": "",
    "name": "foo",
    "value": "",
    "type": "org",
    "images": [
      "alpine"
    ]
  },
  {
    "id": 3,
    "org": "github",
    "repo": "",
    "team": "octokitties",
    "name": "foo",
    "value": "",
    "type": "shared",
    "images": [
      "alpine"
    ]
  }
]`
)

// getSecrets returns mock JSON for a http GET.
func getSecrets(c *gin.Context) {
	data := []byte(SecretsResp)

	var body []api.Secret
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getSecret has a param :name returns mock JSON for a http GET.
//
// Pass "not-found" to :name to test receiving a http 404 response.
func getSecret(c *gin.Context) {
	n := c.Param("name")

	if strings.Contains(n, "not-found") {
		msg := fmt.Sprintf("Secret %s does not exist", n)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	data := []byte(SecretResp)

	var body api.Secret
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addSecret returns mock JSON for a http POST.
func addSecret(c *gin.Context) {
	data := []byte(SecretResp)

	var body api.Secret
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateSecret has a param :name returns mock JSON for a http PUT.
//
// Pass "not-found" to :name to test receiving a http 404 response.
func updateSecret(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		n := c.Param("name")

		if strings.Contains(n, "not-found") {
			msg := fmt.Sprintf("Repo or team %s does not exist for secret", n)

			c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

			return
		}
	}

	data := []byte(SecretResp)

	var body api.Secret
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeSecret has a param :name returns mock JSON for a http DELETE.
//
// Pass "not-found" to :name to test receiving a http 404 response.
func removeSecret(c *gin.Context) {
	n := c.Param("name")

	if strings.Contains(n, "not-found") {
		msg := fmt.Sprintf("Secret %s does not exist", n)

		c.AbortWithStatusJSON(http.StatusNotFound, api.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Secret %s removed", n))
}
