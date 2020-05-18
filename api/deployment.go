// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation POST /api/v1/deployments/{org}/{repo} deployment UpdateDeployment
//
// Create a deployment for the configured backend
//
// ---
// x-success_http_code: '501'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// responses:
//   '501':
//     description: The endpoint is not implemented
//     schema:
//       type: string

// CreateDeployment represents the API handler to
// create a deployment in the configured backend.
func CreateDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// swagger:operation GET /api/v1/deployments/{org}/{repo} deployment GetDeployments
//
// Get a list of deployments for the configured backend
//
// ---
// x-success_http_code: '501'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// responses:
//   '501':
//     description: The endpoint is not implemented
//     schema:
//       type: string

// GetDeployments represents the API handler to capture
// a list of deployments from the configured backend.
func GetDeployments(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// swagger:operation GET /api/v1/deployments/{org}/{repo}/{deployment} deployment GetDeployment
//
// Get a deployment from the configured backend
//
// ---
// x-success_http_code: '501'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: deployment
//   description: Name of the org
//   required: true
//   type: string
// responses:
//   '501':
//     description: The endpoint is not implemented
//     schema:
//       type: string

// GetDeployment represents the API handler to
// capture a deployment from the configured backend.
func GetDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// swagger:operation PUT /api/v1/deployments/{org}/{repo}/{deployment} deployment UpdateDeployment
//
// Update a deployment from the configured backend
//
// ---
// x-success_http_code: '501'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: deployment
//   description: Name of the org
//   required: true
//   type: string
// responses:
//   '501':
//     description: The endpoint is not implemented
//     schema:
//       type: string

// UpdateDeployment represents the API handler to
// update a deployment in the configured backend.
func UpdateDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// swagger:operation DELETE /api/v1/deployments/{org}/{repo}/{deployment} deployment DeleteDeployment
//
// Delete a deployment from the configured backend
//
// ---
// x-success_http_code: '501'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: deployment
//   description: Name of the org
//   required: true
//   type: string
// responses:
//   '501':
//     description: The endpoint is not implemented
//     schema:
//       type: string

// DeleteDeployment represents the API handler to
// remove a deployment from the configured backend.
func DeleteDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}
