// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types"
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

// Login represents the API handler to
// process a user logging in to Vela.
func Login(c *gin.Context) {
	// load the metadata
	m := c.MustGet("metadata").(*types.Metadata)
	// capture an error if present
	err := c.Request.FormValue("error")
	if len(err) > 0 {
		// redirect to initial login screen with error code
		c.Redirect(http.StatusTemporaryRedirect, "/login/error?code="+err)
	}

	// redirect to our authentication handler
	t := c.Request.FormValue("type")
	p := c.Request.FormValue("port")
	r := ""
	path := "/authenticate"

	switch t {
	case "web":
		r = fmt.Sprintf("%s/authenticate/%s", m.Vela.Address, t)
	case "cli":
		if len(p) > 0 {
			r = fmt.Sprintf("%s/authenticate/%s/%s", m.Vela.Address, t, p)
		}
	}

	if len(r) > 0 {
		path += fmt.Sprintf("?redirect_uri=%s", url.QueryEscape(r))
	}

	c.Redirect(http.StatusTemporaryRedirect, path)
}
