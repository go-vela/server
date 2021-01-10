// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

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

	// capture query params
	t := c.Request.FormValue("type")
	p := c.Request.FormValue("port")

	// temp variable to hold redirect destination
	r := ""

	// default path (headless mode)
	path := "/authenticate"

	// handle web and cli logins
	switch t {
	case "web":
		r = fmt.Sprintf("%s/authenticate/%s", m.Vela.Address, t)

		logrus.Debugf("web login request, setting redirect to: %s", r)
	case "cli":
		// port must be supplied
		if len(p) > 0 {
			r = fmt.Sprintf("%s/authenticate/%s/%s", m.Vela.Address, t, p)

			logrus.Debugf("cli login request, setting redirect to: %s", r)
		}

		logrus.Debug("cli login request, but port was not defined")
	}

	// if we a redirecting to non-default destination,
	// prep and append the redirect
	if len(r) > 0 {
		v := &url.Values{}
		v.Add("redirect_uri", r)

		path = fmt.Sprintf("%s?%s", path, v.Encode())
	}

	// redirect to our authentication handler
	c.Redirect(http.StatusTemporaryRedirect, path)
}
