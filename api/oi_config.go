// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/internal/token"
)

// swagger:operation GET /_services/token/.well-known/openid-configuration token GetOpenIDConfig
//
// Get the Vela OIDC service configuration
//
// ---
// produces:
// - application/json
// parameters:
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the Vela OpenID Configuration
//     schema:
//       "$ref": "#/definitions/OpenIDConfig"

// GetOpenIDConfig represents the API handler for requests for configurations in the Vela OpenID service.
func GetOpenIDConfig(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	tm := c.MustGet("token-manager").(*token.Manager)

	l.Debug("reading OpenID configuration")

	config := types.OpenIDConfig{
		Issuer:      tm.Issuer,
		JWKSAddress: fmt.Sprintf("%s/.well-known/jwks", tm.Issuer),
		ClaimsSupported: []string{
			"sub",
			"exp",
			"iat",
			"iss",
			"aud",
			"build_number",
			"build_id",
			"repo",
			"token_type",
			"actor",
			"actor_scm_id",
			"commands",
			"image",
			"request",
			"event",
			"sha",
			"ref",
		},
		ResponseTypesSupported: []string{
			"id_token",
		},
		Algorithms: []string{
			jwt.SigningMethodRS256.Name,
		},
		SubjectTypesSupported: []string{
			"public",
		},
	}

	c.JSON(http.StatusOK, config)
}
