// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

const (
	// TokenRefreshResp represents a JSON return for a token refresh.
	//nolint:gosec // not a hardcoded credential
	TokenRefreshResp = `{
  "token": "header.payload.signature"
}`

	// OpenIDConfigResp represents a JSON return for an OpenID configuration.
	OpenIDConfigResp = `{
  "issuer": "https://vela.com/_services/token",
  "jwks_uri": "https://vela.com/_services/token/.well-known/jwks",
  "supported_claims": [
    "sub",
    "exp",
    "iat",
    "iss",
    "aud",
    "branch",
    "build_number",
    "build_id",
    "repo",
    "token_type",
    "actor",
    "actor_scm_id",
    "commands",
    "image",
    "image_name",
    "image_tag",
    "request"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ]
}`

	// JWKSResp represents a JSON return for the JWKS.
	JWKSResp = `{
  "keys": [
	{
  	"e": "AQAB",
  	"kid": "f7ec4ab7-c9a2-440e-bfb3-83b6599479ea",
  	"kty": "RSA",
  	"n": "weh9G_J6yZEugOFo6MQ057t_ExafteA_zVRS3CEPWiOgBLLRymh-KS6aCW-kHVuyBsnWNrCcc5cRJ6ISFnQMtkJtbpV_72qbw0zhFLiYomZDh5nb5dqCoiWIVNG8_a_My9jhXAIghs8MLbG-_Tj9jZb3K3n3Ies-Cg1E5SWO3YX8I1_X7ZlgqhEbktoy2RvR_crQA_fi1jRW5Q6PldIJmu4FIeXN_ny_sgg6ZQtTImFderUy1aUxUnpjilU-yv13eJejYQnJ7rExJVsDqq3B_CnYD2ioJC6b7aoEPvCpZ_1VgTTnQt6nedmr2Hih3GHgDNsM-BFr63aG3qZ5v9bVRw"
	}
  ]
`
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

// openIDConfig returns a mock response for a http GET.
func openIDConfig(c *gin.Context) {
	data := []byte(OpenIDConfigResp)

	var body api.OpenIDConfig
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getJWKS returns a mock response for a http GET.
func getJWKS(c *gin.Context) {
	data := []byte(JWKSResp)

	var body jwk.RSAPublicKey
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}
