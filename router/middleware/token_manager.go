// SPDX-License-Identifier: Apache-2.0

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
