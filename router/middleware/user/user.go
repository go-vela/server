// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"net/http"
	"strings"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the user in the given context.
func Retrieve(c *gin.Context) *library.User {
	return FromContext(c)
}

// Establish sets the user in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := claims.Retrieve(c)

		// if token is not a user token, establish empty user to better handle nil checks
		if !strings.EqualFold(cl.TokenType, constants.UserAccessTokenType) {
			u := new(library.User)

			ToContext(c, u)
			c.Next()

			return
		}

		logrus.Debugf("parsing user access token")

		// lookup user in claims subject in the database
		u, err := database.FromContext(c).GetUserForName(cl.Subject)
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		ToContext(c, u)
		c.Next()
	}
}
