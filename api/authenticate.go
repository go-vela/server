// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/token"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /authenticate authenticate GetAuthenticate
//
// Start the OAuth flow with the Vela API
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// responses:
//   '200':
//     description: Successfully authenticated
//     schema:
//       type: string
// responses:
//   '307':
//     description: Redirected for authentication
//     schema:
//       type: string
//   '401':
//     description: Unable to authenticate
//     schema:
//       type: string
//   '503':
//     description: Service unavailable
//     schema:
//       type: string

// Authenticate represents the API handler to
// process a user logging in to Vela from
// the API or UI.
func Authenticate(c *gin.Context) {
	var err error

	// capture the OAuth state if present
	oAuthState := c.Request.FormValue("state")

	// capture the OAuth code if present
	code := c.Request.FormValue("code")
	if len(code) == 0 {
		// start the initial OAuth workflow
		oAuthState, err = source.FromContext(c).Login(c.Writer, c.Request)
		if err != nil {
			retErr := fmt.Errorf("unable to login user: %w", err)

			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}

	// complete the OAuth workflow and authenticates the user
	newUser, err := source.FromContext(c).Authenticate(c.Writer, c.Request, oAuthState)
	if err != nil {
		retErr := fmt.Errorf("unable to authenticate user: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// this will happen if the user is redirected by the
	// source provider as part of the authorization workflow.
	if newUser == nil {
		return
	}

	// send API call to capture the user logging in
	u, err := database.FromContext(c).GetUserName(newUser.GetName())
	// create a new user account
	if len(u.GetName()) == 0 || err != nil {
		// create unique id for the user
		uid, err := uuid.NewRandom()
		if err != nil {
			retErr := fmt.Errorf("unable to create UID for user %s: %w", u.GetName(), err)

			util.HandleError(c, http.StatusServiceUnavailable, retErr)

			return
		}

		// create the user account
		u := new(library.User)
		u.SetName(newUser.GetName())
		u.SetToken(newUser.GetToken())
		u.SetHash(
			base64.StdEncoding.EncodeToString(
				[]byte(uid.String()),
			),
		)
		u.SetActive(true)
		u.SetAdmin(false)

		// compose jwt tokens for user
		rt, at, err := token.Compose(c, u)
		if err != nil {
			retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

			util.HandleError(c, http.StatusServiceUnavailable, retErr)

			return
		}

		// store the refresh token with the user object
		u.SetRefreshToken(rt)

		// send API call to create the user in the database
		err = database.FromContext(c).CreateUser(u)
		if err != nil {
			retErr := fmt.Errorf("unable to create user %s: %w", u.GetName(), err)

			util.HandleError(c, http.StatusServiceUnavailable, retErr)

			return
		}

		// return the jwt access token
		c.JSON(http.StatusOK, library.Login{Token: &at})

		return
	}

	// update the user account
	u.SetToken(newUser.GetToken())
	u.SetActive(true)

	// compose jwt tokens for user
	rt, at, err := token.Compose(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	// store the refresh token with the user object
	u.SetRefreshToken(rt)

	// send API call to update the user in the database
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	// return the user with their jwt access token
	c.JSON(http.StatusOK, library.Login{Token: &at})
}

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
	t := c.Param("type")
	p := c.Param("port")

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

// swagger:operation POST /authenticate/token authenticate PostAuthenticate
//
// Authenticate to Vela via personal access token.
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// responses:
//   '200':
//     description: Successfully authenticated
//     schema:
//       type: string
// responses:
//   '401':
//     description: Unable to authenticate
//     schema:
//       type: string
//   '503':
//     description: Service unavailable
//     schema:
//       type: string

// AuthenticateToken represents the API handler to
// process a user logging in using PAT to Vela from
// the API.
func AuthenticateToken(c *gin.Context) {
	// attempt to get user from source
	u, err := source.FromContext(c).AuthenticateToken(c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to authenticate user: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// check if the user exists
	u, err = database.FromContext(c).GetUserName(u.GetName())
	if err != nil {
		retErr := fmt.Errorf("user not found")

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// We don't need refresh token for this scenario
	// We only need access token and are configured based on the config defined
	m := c.MustGet("metadata").(*types.Metadata)
	at, err := token.CreateAccessToken(u, m.Vela.AccessTokenDuration)

	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)
	}

	// return the user with their jwt access token
	c.JSON(http.StatusOK, library.Login{Token: &at})
}
