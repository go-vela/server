// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
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
//       "$ref": "#/definitions/PublicKey"

// GetJWKS represents the API handler for requests to public keys in the Vela OpenID service.
func GetJWKS(c *gin.Context) {
	// retrieve token manager from context
	keys, err := database.FromContext(c).ListKeySets(c)
	if err != nil {
		retErr := fmt.Errorf("unable to get key sets: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
	}

	c.JSON(http.StatusOK, types.JWKS{Keys: keys})
}
