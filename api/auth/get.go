// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation GET /authenticate authenticate GetAuthenticate
//
// Start OAuth flow or exchange tokens
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
// - in: query
//   name: redirect_uri
//   description: the url where the user will be sent after authorization
//   type: string
// responses:
//   '200':
//     description: Successfully authenticated
//     headers:
//       Set-Cookie:
//         type: string
//     schema:
//       "$ref": "#/definitions/Token"
//   '307':
//     description: Redirected for authentication
//   '401':
//     description: Unable to authenticate
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Service unavailable
//     schema:
//       "$ref": "#/definitions/Error"

// Authenticate represents the API handler to
// process a user logging in to Vela from
// the API or UI.
func Authenticate(c *gin.Context) {
	var err error

	tm := c.MustGet("token-manager").(*token.Manager)

	// capture the OAuth state if present
	oAuthState := c.Request.FormValue("state")

	// capture the OAuth code if present
	code := c.Request.FormValue("code")
	if len(code) == 0 {
		// start the initial OAuth workflow
		oAuthState, err = scm.FromContext(c).Login(c.Writer, c.Request)
		if err != nil {
			retErr := fmt.Errorf("unable to login user: %w", err)

			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}

	// complete the OAuth workflow and authenticates the user
	newUser, err := scm.FromContext(c).Authenticate(c.Writer, c.Request, oAuthState)
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
	u, err := database.FromContext(c).GetUserForName(newUser.GetName())
	// create a new user account
	if len(u.GetName()) == 0 || err != nil {
		// create the user account
		u := new(library.User)
		u.SetName(newUser.GetName())
		u.SetToken(newUser.GetToken())
		u.SetActive(true)
		u.SetAdmin(false)

		// compose jwt tokens for user
		rt, at, err := tm.Compose(c, u)
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
		c.JSON(http.StatusOK, library.Token{Token: &at})

		return
	}

	// update the user account
	u.SetToken(newUser.GetToken())
	u.SetActive(true)

	// compose jwt tokens for user
	rt, at, err := tm.Compose(c, u)
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
	c.JSON(http.StatusOK, library.Token{Token: &at})
}
