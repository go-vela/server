// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /login authenticate GetLogin
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

// swagger:operation GET /logout authenticate Logout
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

// swagger:operation POST /login authenticate PostLogin
//
// Login to the Vela api
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Login payload that we expect from the user
//   required: true
//   schema:
//     "$ref": "#/definitions/Login"
// responses:
//   '200':
//     description: Successful login to the Vela API
//     schema:
//       "$ref": "#/definitions/Login"
//   '400':
//     description: Unable to login to the Vela API
//     schema:
//       type: string
//   '401':
//     description: Unable to login to the Vela API
//     schema:
//       type: string
//   '503':
//     description: Unable to login to the Vela API
//     schema:
//       type: string

// Login represents the API handler to
// process a user logging in to Vela.
func Login(c *gin.Context) {
	// check if request was a POST
	if strings.EqualFold(c.Request.Method, "POST") {
		// assume all POST requests are coming from the CLI
		AuthenticateCLI(c)

		return
	}

	// capture an error if present
	err := c.Request.FormValue("error")
	if len(err) > 0 {
		// redirect to initial login screen with error code
		c.Redirect(http.StatusTemporaryRedirect, "/login/error?code="+err)
	}

	// redirect to our authentication handler
	c.Redirect(http.StatusTemporaryRedirect, "/authenticate")
}
