// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// Allowlist is a middleware function that attaches the allowlist used
// to limit which repos can be activated within the system.
func Allowlist(allowlist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("allowlist", allowlist)
		c.Next()
	}
}
