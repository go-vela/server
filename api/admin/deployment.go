// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /api/v1/admin/deployments admin AllDeployments
//
// Get all of the deployments in the database
//
// ---
// x-success_http_code: '501'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// responses:
//   '501':
//     description: This endpoint is not implemented
//     schema:
//       type: string

// AllDeployments represents the API handler to
// captures all deployments stored in the database.
func AllDeployments(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// swagger:operation PUT /api/v1/admin/deployment admin UpdateDeployment
//
// Get All
//
// ---
// x-success_http_code: '501'
// x-incident_priority: P4
// produces:
// - application/json
// parameters:
// responses:
//   '501':
//     description: This endpoint is not implemented
//     schema:
//       type: string

// UpdateDeployment represents the API handler to
// update any deployment stored in the database.
func UpdateDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}
