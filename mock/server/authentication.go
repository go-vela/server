// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

const (
	// TokenRefreshResp represents a JSON return for a token refresh.
	//nolint:gosec // not a hardcoded credential
	TokenRefreshResp = `{
  "token": "header.payload.signature"
}`
)

// getTokenRefresh returns mock JSON for a http GET.
func getTokenRefresh(c *gin.Context) {
	data := []byte(TokenRefreshResp)

	var body library.Token
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getAuthenticate returns mock response for a http GET.
//
// Don't pass "state" and "code" params to receive an error response.
func getAuthenticate(c *gin.Context) {
	data := []byte(TokenRefreshResp)

	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	err := "error"

	if len(state) == 0 && len(code) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &err})

		return
	}

	var body library.Token
	_ = json.Unmarshal(data, &body)

	c.SetCookie(constants.RefreshTokenName, "refresh", 2, "/", "", true, true)

	c.JSON(http.StatusOK, body)
}

// getAuthenticateFromToken returns mock response for a http POST.
//
// Don't pass "Token" in header to receive an error message.
func getAuthenticateFromToken(c *gin.Context) {
	data := []byte(TokenRefreshResp)
	err := "error"

	token := c.Request.Header.Get("Token")
	if len(token) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &err})
	}

	var body library.Token
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// validateToken returns mock response for a http GET.
//
// Don't pass "Authorization" in header to receive an unauthorized error message.
func validateToken(c *gin.Context) {
	err := "error"

	token := c.Request.Header.Get("Authorization")
	if len(token) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &err})
	}

	c.JSON(http.StatusOK, "vela-server")
}

// validateOAuthToken returns mock response for a http GET.
//
// Don't pass "Authorization" in header to receive an unauthorized error message.
func validateOAuthToken(c *gin.Context) {
	err := "error"

	token := c.Request.Header.Get("Authorization")
	if len(token) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Error{Message: &err})
	}

	c.JSON(http.StatusOK, "oauth token was created by vela")
}
