// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /logout authenticate GetLogout
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
//   '503':
//     description: Service unavailable
//     schema:
//       type: string

// Logout represents the API handler to
// process a user logging out of Vela.
// Primarily, it deletes the current
// refresh token cookie.
func Logout(c *gin.Context) {
	// grab the metadata to help deal with the cookie
	m := c.MustGet("metadata").(*types.Metadata)
	u := user.Retrieve(c)

	logrus.Infof("logging out user: %s", u.GetName())

	// parse the address for the backend server
	// so we can set it for the cookie domain
	addr, err := url.Parse(m.Vela.Address)
	if err != nil {
		// silently fail
		logrus.Error("unable to parse Vela address during logout")
	}

	// remove the refresh token from the cookies, Max-Age value -1 will do it
	c.SetCookie(constants.RefreshTokenName, "", -1, "/", addr.Hostname(), c.Value("securecookie").(bool), true)

	// unset the refresh token for the user
	u.SetRefreshToken("")

	// send API call to update the user in the database
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	// return 200 for successful logout
	c.JSON(http.StatusOK, "ok")
}
