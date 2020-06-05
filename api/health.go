// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /health router Health
//
// Check if the Vela API is available
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successfully 'ping'-ed Vela API
//     schema:
//       type: string

// Health represents the API handler to
// report the health status for Vela.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
