// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/internal/token"
)

// TokenManager is a middleware function that attaches the token manager
// to the context of every http.Request.
func TokenManager(m *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("token-manager", m)
		c.Next()
	}
}
