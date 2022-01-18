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
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

// nolint:gosec // these are mock responses
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
  "events": [
    "push"
  ]
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
    ],
    "events": [
      "push"
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
    ],
    "events": [
      "push"
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
    ],
    "events": [
      "push"
    ]
  }
]`
)

// getSecrets returns mock JSON for a http GET.
func getSecrets(c *gin.Context) {
	data := []byte(SecretsResp)

	var body []library.Secret
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

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(SecretResp)

	var body library.Secret
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addSecret returns mock JSON for a http POST.
func addSecret(c *gin.Context) {
	data := []byte(SecretResp)

	var body library.Secret
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

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(SecretResp)

	var body library.Secret
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

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Secret %s removed", n))
}
