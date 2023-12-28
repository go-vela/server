// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// DefaultBuildLimit is a middleware function that attaches the defaultLimit
// to enable the server to override the default build limit.
func DefaultBuildLimit(defaultBuildLimit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("defaultBuildLimit", defaultBuildLimit)
		c.Next()
	}
}
