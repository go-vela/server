// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/version"
)

// swagger:operation GET /version router Version
//
// Get the version of the Vela API
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successfully retrieved the Vela API version
//     schema:
//       type: string

// Version represents the API handler to
// report the version number for Vela.
func Version(c *gin.Context) {
	c.JSON(http.StatusOK, version.New())
}
