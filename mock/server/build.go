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
	// BuildResp represents a JSON return for a single build.
	BuildResp = `{
  "id": 1,
  "repo_id": 1,
  "number": 1,
  "parent": 1,
  "event": "push",
  "status": "created",
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
  "branch": "master",
  "ref": "refs/heads/master",
  "base_ref": "",
  "host": "example.company.com",
  "runtime": "docker",
  "distribution": "linux"
}`

	// BuildsResp represents a JSON return for one to many builds.
	BuildsResp = `[
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
    "branch": "master",
    "ref": "refs/heads/master",
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
    "branch": "master",
    "ref": "refs/heads/master",
    "base_ref": "",
    "host": "82823eb770b0",
    "runtime": "",
    "distribution": ""
  }
]`

	// BuildLogsResp represents a JSON return for build logs.
	BuildLogsResp = `[
  {
    "id": 1,
    "step_id": 1,
    "build_id": 1,
    "repo_id": 1,
    "data": "SGVsbG8sIFdvcmxkIQ=="
  },
  {
    "id": 2,
    "step_id": 2,
    "build_id": 1,
    "repo_id": 1,
    "data": "SGVsbG8sIFdvcmxkIQ=="
  }
]`

	// BuildQueueResp represents a JSON return for build queue.
	BuildQueueResp = `[
  {
    "status": "running",
    "created": 1616467142,
    "number": 6,
    "full_name": "github/octocat"
  },
  {
    "status": "pending",
    "created": 1616467142,
    "number": 7,
    "full_name": "github/octocat"
  }
]`

	// BuildTokenResp represents a JSON return for requesting a build token
	//
	//nolint:gosec // not actual credentials
	BuildTokenResp = `{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJidWlsZF9pZCI6MSwicmVwbyI6ImZvby9iYXIiLCJzdWIiOiJPY3RvY2F0IiwiaWF0IjoxNTE2MjM5MDIyfQ.hD7gXpaf9acnLBdOBa4GOEa5KZxdzd0ZvK6fGwaN4bc"
  }`

	// BuildItineraryResp represents a JSON return for requesting a build itinerary.
	BuildItineraryResp = `{
    "id": 1
    "build_id": 1,
    "data": "eyAKICAgICJpZCI6ICJzdGVwX25hbWUiLAogICAgInZlcnNpb24iOiAiMSIsCiAgICAibWV0YWRhdGEiOnsKICAgICAgICAiY2xvbmUiOnRydWUsCiAgICAgICAgImVudmlyb25tZW50IjpbInN0ZXBzIiwic2VydmljZXMiLCJzZWNyZXRzIl19LAogICAgIndvcmtlciI6e30sCiAgICAic3RlcHMiOlsKICAgICAgICB7CiAgICAgICAgICAgICJpZCI6InN0ZXBfZ2l0aHViX29jdG9jYXRfMV9pbml0IiwKICAgICAgICAgICAgImRpcmVjdG9yeSI6Ii92ZWxhL3NyYy9naXRodWIuY29tL2dpdGh1Yi9vY3RvY2F0IiwKICAgICAgICAgICAgImVudmlyb25tZW50IjogeyJCVUlMRF9BVVRIT1IiOiJPY3RvY2F0In0KICAgICAgICB9CiAgICBdCn0KCg=="
  }`
)

// getBuilds returns mock JSON for a http GET.
func getBuilds(c *gin.Context) {
	data := []byte(BuildsResp)

	var body []library.Build
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getBuild has a param :build returns mock JSON for a http GET.
func getBuild(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(BuildResp)

	var body library.Build
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getLogs has a param :build returns mock JSON for a http GET.
//
// Pass "0" to :build to test receiving a http 404 response.
func getLogs(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(BuildLogsResp)

	var body []library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addBuild returns mock JSON for a http POST.
func addBuild(c *gin.Context) {
	data := []byte(BuildResp)

	var body library.Build
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateBuild has a param :build returns mock JSON for a http PUT.
//
// Pass "0" to :build to test receiving a http 404 response.
func updateBuild(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		b := c.Param("build")

		if strings.EqualFold(b, "0") {
			msg := fmt.Sprintf("Build %s does not exist", b)

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(BuildResp)

	var body library.Build
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeBuild has a param :build returns mock JSON for a http DELETE.
//
// Pass "0" to :build to test receiving a http 404 response.
func removeBuild(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Build %s removed", b))
}

// restartBuild has a param :build returns mock JSON for a http POST.
//
// Pass "0" to :build to test receiving a http 404 response.
func restartBuild(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(BuildResp)

	var body library.Build
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// cancelBuild has a param :build returns mock JSON for a http DELETE.
//
// Pass "0" to :build to test receiving a http 404 response.
func cancelBuild(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("Build %s does not exist", b)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, BuildResp)
}

// buildQueue has a param :after returns mock JSON for a http GET.
//
// Pass "0" to :after to test receiving a http 200 response with no builds.
func buildQueue(c *gin.Context) {
	b := c.Param("after")

	if strings.EqualFold(b, "0") {
		c.AbortWithStatusJSON(http.StatusOK, []string{})

		return
	}

	data := []byte(BuildQueueResp)

	var body []library.BuildQueue
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// buildToken has a param :build returns mock JSON for a http GET.
//
// Pass "0" to :build to test receiving a http 404 response. Pass "2"
// to :build to test receiving a http 400 response.
func buildToken(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		c.AbortWithStatusJSON(http.StatusNotFound, "")

		return
	}

	if strings.EqualFold(b, "2") {
		c.AbortWithStatusJSON(http.StatusBadRequest, "")
	}

	data := []byte(BuildTokenResp)

	var body library.Token
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// buildItinerary has a param :build returns mock JSON for a http GET.
//
// Pass "0" to :build to test receiving a http 500 response.
func buildItinerary(c *gin.Context) {
	b := c.Param("build")

	if strings.EqualFold(b, "0") {
		msg := fmt.Sprintf("unable to get build itinerary for build %s", b)

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.Error{Message: &msg})

		return
	}

	data := []byte(BuildItineraryResp)

	var body library.BuildItinerary
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
