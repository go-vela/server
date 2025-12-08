// SPDX-License-Identifier: Apache-2.0

package claims

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/auth"
	"github.com/go-vela/server/util"
)

// Retrieve gets the claims in the given context.
func Retrieve(c *gin.Context) *token.Claims {
	return FromContext(c)
}

// Establish sets the claims in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
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

		// if this is an installation token, no claims
		if strings.HasPrefix(at, "ghs_") {
			installToken, err := cache.FromContext(c).GetInstallToken(c.Request.Context(), at)
			if err != nil || installToken == nil {
				retErr := fmt.Errorf("unable to validate installation token: %w", err)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return
			}

			c.Set("app-installation-token", installToken)
			c.Next()

			return
		}

		// parse and validate the token and return the associated the user
		claims, err = tm.ParseToken(at)
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		l = l.WithFields(logrus.Fields{
			"claim_subject": claims.Subject,
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, claims)
		c.Next()
	}
}
