// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
