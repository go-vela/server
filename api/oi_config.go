// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/internal"
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
	m := c.MustGet("metadata").(*internal.Metadata)
	config := types.OpenIDConfig{
		Issuer:      fmt.Sprintf("%s/_services/token", m.Vela.Address),
		JWKSAddress: fmt.Sprintf("%s/%s", m.Vela.Address, "_services/token/.well-known/jwks"),
		SupportedClaims: []string{
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
		Algorithms: []string{
			jwt.SigningMethodRS256.Name,
		},
	}

	c.JSON(http.StatusOK, config)
}
