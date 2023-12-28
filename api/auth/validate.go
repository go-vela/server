// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /validate-token authenticate ValidateServerToken
//
// Validate a server token
//
// ---
// produces:
// - application/json
// security:
//   - CookieAuth: []
// responses:
//   '200':
//     description: Successfully validated a token
//     schema:
//       type: string
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// ValidateServerToken will validate if a token was issued by the server
// if it is provided in the auth header.
func ValidateServerToken(c *gin.Context) {
	cl := claims.Retrieve(c)

	if !strings.EqualFold(cl.Subject, "vela-server") {
		retErr := fmt.Errorf("token is not a valid server token")

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, "valid server token")
}
