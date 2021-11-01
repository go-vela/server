// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	// DeploymentResp represents a JSON return for a single build.
	DeploymentResp = `{
  "id": 1,
  "repo_id": 1,
  "url": "https://api.github.com/repos/github/octocat/deployments/1",
  "user": "octocat",
  "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
  "ref": "master",
  "task": "deploy:vela",
  "target": "production",
  "description": "Deployment request from Vela"
}`

	// DeploymentsResp represents a JSON return for one to many builds.
	DeploymentsResp = `[
  {
    "id": 2,
    "repo_id": 1,
    "url": "https://api.github.com/repos/github/octocat/deployments/2",
    "user": "octocat",
    "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
    "ref": "master",
    "task": "deploy:vela",
    "target": "production",
    "description": "Deployment request from Vela"
  },
  {
    "id": 1,
    "repo_id": 1,
    "url": "https://api.github.com/repos/github/octocat/deployments/1",
    "user": "octocat",
    "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
    "ref": "master",
    "task": "deploy:vela",
    "target": "production",
    "description": "Deployment request from Vela"
  }
]`
)

// getDeployments returns mock JSON for a http GET.
func getDeployments(c *gin.Context) {
	data := []byte(DeploymentsResp)

	var body []library.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getDeployment has a param :deployment returns mock JSON for a http GET.
func getDeployment(c *gin.Context) {
	d := c.Param("deployment")

	if strings.EqualFold(d, "0") {
		msg := fmt.Sprintf("Deployment %s does not exist", d)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(DeploymentResp)

	var body library.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addDeployment returns mock JSON for a http POST.
func addDeployment(c *gin.Context) {
	data := []byte(DeploymentResp)

	var body library.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateDeployment returns mock JSON for a http PUT.
func updateDeployment(c *gin.Context) {
	data := []byte(DeploymentResp)

	var body library.Deployment
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
