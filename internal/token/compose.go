// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// Compose generates a refresh and access token pair unique
// to the provided user and sets a secure cookie.
// It uses the user's hash to sign the token. to
// guarantee the signature is unique per token. The refresh
// token is returned to store with the user
// in the database.
func (tm *Manager) Compose(c *gin.Context, u *library.User) (string, string, error) {
	// grab the metadata from the context to pull in provided
	// cookie duration information
	m := c.MustGet("metadata").(*types.Metadata)

	// mint token options for refresh token
	rmto := MintTokenOpts{
		User:          u,
		TokenType:     constants.UserRefreshTokenType,
		TokenDuration: tm.UserRefreshTokenDuration,
	}

	// create a refresh token with the provided options
	refreshToken, err := tm.MintToken(&rmto)
	if err != nil {
		return "", "", err
	}

	// mint token options for access token
	amto := MintTokenOpts{
		User:          u,
		TokenType:     constants.UserAccessTokenType,
		TokenDuration: tm.UserAccessTokenDuration,
	}

	// create an access token with the provided options
	accessToken, err := tm.MintToken(&amto)
	if err != nil {
		return "", "", err
	}

	// parse the address for the backend server
	// so we can set it for the cookie domain
	addr, err := url.Parse(m.Vela.Address)
	if err != nil {
		return "", "", err
	}

	refreshExpiry := int(tm.UserRefreshTokenDuration.Seconds())

	// set the SameSite value for the cookie
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#samesite-attribute
	// We set to Lax because we will have links from source provider web UI.
	// Setting this to Strict would force a login when navigating via source provider web UI links.
	c.SetSameSite(http.SameSiteLaxMode)
	// set the cookie with the refresh token as a HttpOnly, Secure cookie
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#httponly-attribute
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#secure-attribute
	c.SetCookie(constants.RefreshTokenName, refreshToken, refreshExpiry, "/", addr.Hostname(), c.Value("securecookie").(bool), true)

	// return the refresh and access tokens
	return refreshToken, accessToken, nil
}
