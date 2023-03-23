// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/auth"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, library.Token{Token: &newAccessToken})
}

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
//       "$ref": "#/definitions/Claims"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// ValidateServerToken will return the claims of a valid server token
// if it is provided in the auth header.
func ValidateServerToken(c *gin.Context) {
	cl := claims.Retrieve(c)

	if !strings.EqualFold(cl.Subject, "vela-server") {
		retErr := fmt.Errorf("token is not a valid server token")

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, cl)
}
