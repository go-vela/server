// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/settings"
)

// Allowlist is a middleware function that attaches the allowlist used
// to limit which repos can be activated within the system.
func Allowlist() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := settings.FromContext(c)
		c.Set("allowlist", s.GetRepoAllowlist())
		c.Next()
	}
}
