// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
)

// swagger:operation GET /logout router GetLogout
//
// Log out of the Vela api
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successfully logged out
//     schema:
//       type: string

// Logout represents the API handler to
// process a user logging out of Vela.
// Primarily, it deletes the current
// refresh token cookie.
func Logout(c *gin.Context) {
	// grab the metadata to help deal with the cookie
	m := c.MustGet("metadata").(*types.Metadata)

	// parse the address for the backend server
	// so we can set it for the cookie domain
	addr, err := url.Parse(m.Vela.Address)
	if err != nil {
		// silently fail
		// TODO: reconsider?
		return
	}

	// remove the refresh token from the cookies, Max-Age value -1 will do it
	c.SetCookie(constants.RefreshTokenName, "", -1, "/", addr.Hostname(), false, true)

	// return 200 for successful logout
	c.JSON(http.StatusOK, "ok")
}
