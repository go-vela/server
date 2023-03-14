// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package claims

import (
	"net/http"
	"strings"

	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/auth"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"

	"github.com/gin-gonic/gin"
)

// Retrieve gets the claims in the given context.
func Retrieve(c *gin.Context) *token.Claims {
	return FromContext(c)
}

// Establish sets the claims in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		tm := c.MustGet("token-manager").(*token.Manager)
		// get the access token from the request
		at, err := auth.RetrieveAccessToken(c.Request)
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		claims := new(token.Claims)

		// special handling for workers if symmetric token is provided
		if secret, ok := c.Value("secret").(string); ok {
			if strings.EqualFold(at, secret) {
				claims.Subject = "vela-worker"
				claims.TokenType = constants.ServerWorkerTokenType
				ToContext(c, claims)
				c.Next()

				return
			}
		}

		// parse and validate the token and return the associated the user
		claims, err = tm.ParseToken(at)
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		ToContext(c, claims)
		c.Next()
	}
}
