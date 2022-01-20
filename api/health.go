// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /health base Health
//
// Check if the Vela API is available
//
// ---
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
