// SPDX-License-Identifier: Apache-2.0

package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/util"
)

// Retrieve gets the user in the given context.
func Retrieve(c *gin.Context) *api.User {
	return FromContext(c)
}

// Establish sets the user in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)
		ctx := c.Request.Context()

		// if token is not a user token or claims were not retrieved, establish empty user to better handle nil checks
		if cl == nil || !strings.EqualFold(cl.TokenType, constants.UserAccessTokenType) {
			u := new(api.User)

			ToContext(c, u)
			c.Next()

			return
		}

		l.Debugf("parsing user access token")

		// lookup user in claims subject in the database
		u, err := database.FromContext(c).GetUserForName(ctx, cl.Subject)
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		l = l.WithFields(logrus.Fields{
			"user":    u.GetName(),
			"user_id": u.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, u)
		c.Next()
	}
}
