// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /_services/token/.well-known/jwks token GetJWKS
//
// Get the JWKS for the Vela OIDC service
//
// ---
// produces:
// - application/json
// parameters:
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the Vela JWKS
//     schema:
//       "$ref": "#/definitions/JWKSet"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetJWKS represents the API handler for requests to public keys in the Vela OpenID service.
func GetJWKS(c *gin.Context) {
	// retrieve JWKs from the database
	keys, err := database.FromContext(c).ListJWKs(c)
	if err != nil {
		retErr := fmt.Errorf("unable to get key set: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, keys)
}
