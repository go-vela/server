// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"fmt"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/token"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"
	"net/http"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// swagger:operation POST /authenticate authenticate PostAuthenticate
//
// Complete the OAuth flow with the Vela API
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

	c.JSON(http.StatusOK, performUserOperation(c, newUser.GetName(), newUser.GetToken()))
}

// AuthenticateToken represents the API handler to
// process a user logging in using PAT to Vela from
// the API
func AuthenticateToken(c *gin.Context) {
	newUser, err := source.FromContext(c).AuthenticateToken(c.Writer, c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to authenticate user: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, performUserOperation(c, newUser.GetName(), newUser.GetToken()))
}

// AuthenticateCLI represents the API handler to
// process a user logging in to Vela from
// the CLI.
func AuthenticateCLI(c *gin.Context) {
	// capture body from API request
	input := new(library.Login)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// register user with OAuth application
	newUser, err := source.FromContext(c).LoginCLI(input.GetUsername(), input.GetPassword(), input.GetOTP())
	if err != nil {
		retErr := fmt.Errorf("unable to login user: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, performUserOperation(c, newUser.GetName(), newUser.GetToken()))
}

// performUserOperation add or update user in database
// Also, return the login user as response
func performUserOperation(c *gin.Context, userName, authToken string) *library.Login {
	u, err := database.FromContext(c).GetUserName(userName)
	if len(u.GetName()) == 0 || err != nil {
		uid, err := uuid.NewRandom()
		if err != nil {
			retErr := fmt.Errorf("unable to create UID: %v", err)
			util.HandleError(c, http.StatusServiceUnavailable, retErr)
			return nil
		}
		u := new(library.User)
		u.SetName(userName)
		u.SetToken(authToken)
		u.SetHash(
			base64.StdEncoding.EncodeToString(
				[]byte(uid.String()),
			),
		)
		u.SetActive(true)
		u.SetAdmin(false)

		err = database.FromContext(c).CreateUser(u)

		if err != nil {
			retErr := fmt.Errorf("unable to create user %s: %w", u.GetName(), err)

			util.HandleError(c, http.StatusServiceUnavailable, retErr)

			return nil
		}

		// compose JWT token for user
		t, err := token.Compose(u)
		if err != nil {
			retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

			util.HandleError(c, http.StatusServiceUnavailable, retErr)

			return nil
		}
		return &library.Login{Username: u.Name, Token: &t}
	}

	u.SetToken(authToken)
	u.SetActive(true)

	// send API call to update the user in the database
	err = database.FromContext(c).UpdateUser(u)

	// send API call to update the user in the database
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return nil
	}

	// compose JWT token for user
	t, err := token.Compose(u)

	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return nil
	}

	return &library.Login{Username: u.Name, Token: &t}
}
