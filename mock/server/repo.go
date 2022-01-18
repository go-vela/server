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

const (
	// RepoResp represents a JSON return for a single repo.
	RepoResp = `{
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

	// ReposResp represents a JSON return for one to many repos.
	ReposResp = `[
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

// getRepos returns mock JSON for a http GET.
func getRepos(c *gin.Context) {
	data := []byte(ReposResp)

	var body []library.Repo
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getRepo has a param :repo returns mock JSON for a http GET.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func getRepo(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(RepoResp)

	var body library.Repo
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addRepo returns mock JSON for a http POST.
func addRepo(c *gin.Context) {
	data := []byte(RepoResp)

	var body library.Repo
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateRepo has a param :repo returns mock JSON for a http PUT.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func updateRepo(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		r := c.Param("repo")

		if strings.Contains(r, "not-found") {
			msg := fmt.Sprintf("Repo %s does not exist", r)

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(RepoResp)

	var body library.Repo
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeRepo has a param :repo returns mock JSON for a http DELETE.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func removeRepo(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Repo %s removed", r))
}

// repairRepo has a param :repo returns mock JSON for a http PATCH.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func repairRepo(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Repo %s repaired", r))
}

// chownRepo has a param :repo returns mock JSON for a http PATCH.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func chownRepo(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Repo %s changed org", r))
}
