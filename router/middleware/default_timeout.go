// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// DefaultTimeout is a middleware function that attaches the defaultTimeout
// to enable the server to override the default build timeout.
func DefaultTimeout(defaultTimeout int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("defaultTimeout", defaultTimeout)
		c.Next()
	}
}
