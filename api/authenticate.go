// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// swagger:operation GET /authenticate authenticate Authenticate
//
// Authenticate with the Vela API
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing login information
//   required: true
//   schema:
//     "$ref": "#/definitions/Login"
// responses:
//   '200':
//     description: Successfully authenticated
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

		// compose JWT token for user
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

		// return the jwt token
		c.JSON(http.StatusOK, library.Login{Token: &at})

		return
	}

	// update the user account
	u.SetToken(newUser.GetToken())
	u.SetActive(true)

	// compose JWT token for user
	rt, at, err := token.Compose(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	u.SetRefreshToken(rt)

	// send API call to update the user in the database
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	// return the user with their JWT token
	c.JSON(http.StatusOK, library.Login{Token: &at})
}
