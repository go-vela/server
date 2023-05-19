// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /authenticate/web authenticate GetAuthenticateTypeWeb
//
// Authentication entrypoint that builds the right post-auth
// redirect URL for web authentication requests
// and redirects to /authenticate after
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: code
//   description: the code received after identity confirmation
//   type: string
// - in: query
//   name: state
//   description: a random string
//   type: string
// responses:
//   '307':
//     description: Redirected for authentication

// swagger:operation GET /authenticate/cli/{port} authenticate GetAuthenticateTypeCLI
//
// Authentication entrypoint that builds the right post-auth
// redirect URL for CLI authentication requests
// and redirects to /authenticate after
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: port
//   required: true
//   description: the port number
//   type: integer
// - in: query
//   name: code
//   description: the code received after identity confirmation
//   type: string
// - in: query
//   name: state
//   description: a random string
//   type: string
// responses:
//   '307':
//     description: Redirected for authentication

// AuthenticateType handles cases where the OAuth callback was
// overridden by supplying a redirect_uri in the login process.
// It will send the user to the destination to handle the last leg
// in the auth flow - exchanging "code" and "state" for a token.
// This will only handle non-headless flows (ie. web or cli).
func AuthenticateType(c *gin.Context) {
	// load the metadata
	m := c.MustGet("metadata").(*types.Metadata)

	logrus.Info("redirecting for final auth flow destination")

	// capture the path elements
	t := util.PathParameter(c, "type")
	p := util.PathParameter(c, "port")

	// capture the current query parameters -
	// they should contain the "code" and "state" values
	q := c.Request.URL.Query()

	// default redirect location if a user ended up here
	// by providing an unsupported type
	r := fmt.Sprintf("%s/authenticate", m.Vela.Address)

	switch t {
	// cli auth flow
	case "cli":
		r = fmt.Sprintf("http://127.0.0.1:%s", p)
	// web auth flow
	case "web":
		r = fmt.Sprintf("%s%s", m.Vela.WebAddress, m.Vela.WebOauthCallbackPath)
	}

	// append the code and state values
	r = fmt.Sprintf("%s?%s", r, q.Encode())

	c.Redirect(http.StatusTemporaryRedirect, r)
}
