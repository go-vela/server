// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/auth"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /token-refresh authenticate GetRefreshAccessToken
//
// Refresh an access token
//
// ---
// produces:
// - application/json
// security:
//   - CookieAuth: []
// responses:
//   '200':
//     description: Successfully refreshed a token
//     schema:
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// RefreshAccessToken will return a new access token if the provided
// refresh token via cookie is valid.
func RefreshAccessToken(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)

	l.Info("refreshing access token")

	// capture the refresh token
	// TODO: move this into token package and do it internally
	// since we are already passsing context
	rt, err := auth.RetrieveRefreshToken(c.Request)
	if err != nil {
		retErr := fmt.Errorf("refresh token error: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	tm := c.MustGet("token-manager").(*token.Manager)

	// validate the refresh token and return a new access token
	newAccessToken, err := tm.Refresh(c, rt)
	if err != nil {
		retErr := fmt.Errorf("unable to refresh token: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, types.Token{Token: &newAccessToken})
}
