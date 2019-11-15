// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"net/http"
	"strings"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/token"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the user in the given context
func Retrieve(c *gin.Context) *library.User {
	return FromContext(c)
}

// Establish sets the user in the given context
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := token.Retrieve(c.Request)
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		secret := c.MustGet("secret").(string)
		if strings.EqualFold(t, secret) {
			u := new(library.User)
			u.SetName("vela-worker")
			u.SetActive(true)
			u.SetAdmin(true)

			ToContext(c, u)
			c.Next()
			return
		}

		logrus.Debugf("Parsing user token")
		u, err := token.Parse(t, database.FromContext(c))
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		ToContext(c, u)
		c.Next()
	}
}
