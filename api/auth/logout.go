// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /logout authenticate GetLogout
//
// Log out of the Vela api
//
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Successfully logged out
//     schema:
//       type: string
//   '503':
//     description: Logout did not succeed
//     schema:
//       "$ref": "#/definitions/Error"

// Logout represents the API handler to
// process a user logging out of Vela.
// Primarily, it deletes the current
// refresh token cookie.
func Logout(c *gin.Context) {
	// grab the metadata to help deal with the cookie
	m := c.MustGet("metadata").(*types.Metadata)
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	})

	logger.Infof("logging out user %s", u.GetName())

	// parse the address for the backend server
	// so we can set it for the cookie domain
	addr, err := url.Parse(m.Vela.Address)
	if err != nil {
		// silently fail
		logger.Error("unable to parse Vela address during logout")
	}

	// set the same samesite attribute we used to create the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	// remove the refresh token from the cookies, Max-Age value -1 will do it
	c.SetCookie(
		constants.RefreshTokenName, "", -1, "/", addr.Hostname(), c.Value("securecookie").(bool), true,
	)

	// unset the refresh token for the user
	u.SetRefreshToken("")

	// send API call to update the user in the database
	_, err = database.FromContext(c).UpdateUser(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	// return 200 for successful logout
	c.JSON(http.StatusOK, "ok")
}
