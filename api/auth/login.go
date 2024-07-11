// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /login authenticate GetLogin
//
// Log into the Vela API
//
// ---
// parameters:
// - in: query
//   name: type
//   description: The login type ("cli" or "web")
//   type: string
//   enum:
//     - web
//     - cli
// - in: query
//   name: port
//   description: The port number when type=cli
//   type: integer
// responses:
//   '307':
//     description: Redirected to /authenticate

// Login represents the API handler to
// process a user logging in to Vela.
func Login(c *gin.Context) {
	// load the metadata
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)

	// capture query params
	t := util.FormParameter(c, "type")
	p := util.FormParameter(c, "port")

	// temp variable to hold redirect destination
	r := ""

	// default path (headless mode)
	path := "/authenticate"

	l.Info("logging in user")

	// handle web and cli logins
	switch t {
	case "web":
		r = fmt.Sprintf("%s/authenticate/%s", m.Vela.Address, t)

		l.Infof("web login request, setting redirect to: %s", r)
	case "cli":
		// port must be supplied
		if len(p) > 0 {
			r = fmt.Sprintf("%s/authenticate/%s/%s", m.Vela.Address, t, p)

			l.Infof("cli login request, setting redirect to: %s", r)
		}
	}

	// if we a redirecting to non-default destination,
	// prep and append the redirect
	if len(r) > 0 {
		v := &url.Values{}
		v.Add("redirect_uri", r)

		path = fmt.Sprintf("%s?%s", path, v.Encode())
	}

	// redirect to our authentication handler
	// will be either <vela server>/authenticate (headless)
	// or <vela server>/authenticate?redirect_uri=<redirect> (web or cli)
	c.Redirect(http.StatusTemporaryRedirect, path)
}
