// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /login router GetLogin
//
// Log into the Vela api
//
// ---
// x-success_http_code: '307'
// produces:
// - application/json
// parameters:
// responses:
//   '307':
//     description: Redirected to /authenticate
//     schema:
//       type: string

// Login represents the API handler to
// process a user logging in to Vela.
func Login(c *gin.Context) {
	// capture an error if present
	err := c.Request.FormValue("error")
	if len(err) > 0 {
		// redirect to initial login screen with error code
		c.Redirect(http.StatusTemporaryRedirect, "/login/error?code="+err)
	}

	// redirect to our authentication handler
	c.Redirect(http.StatusTemporaryRedirect, "/authenticate")
}
