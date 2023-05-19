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
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// swagger:operation POST /authenticate/token authenticate PostAuthenticateToken
//
// Authenticate to Vela via personal access token
//
// ---
// produces:
// - application/json
// parameters:
// - in: header
//   name: Token
//   type: string
//   required: true
//   description: >
//     scopes: repo, repo:status, user:email, read:user, and read:org
// responses:
//   '200':
//     description: Successfully authenticated
//     schema:
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unable to authenticate
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Service unavailable
//     schema:
//       "$ref": "#/definitions/Error"

// AuthenticateToken represents the API handler to
// process a user logging in using PAT to Vela from
// the API.
func AuthenticateToken(c *gin.Context) {
	// attempt to get user from source
	u, err := scm.FromContext(c).AuthenticateToken(c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to authenticate user: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// check if the user exists
	u, err = database.FromContext(c).GetUserForName(u.GetName())
	if err != nil {
		retErr := fmt.Errorf("user %s not found", u.GetName())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// We don't need refresh token for this scenario
	// We only need access token and are configured based on the config defined
	tm := c.MustGet("token-manager").(*token.Manager)

	// mint token options for access token
	amto := &token.MintTokenOpts{
		User:          u,
		TokenType:     constants.UserAccessTokenType,
		TokenDuration: tm.UserAccessTokenDuration,
	}
	at, err := tm.MintToken(amto)

	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)
	}

	// return the user with their jwt access token
	c.JSON(http.StatusOK, library.Token{Token: &at})
}
