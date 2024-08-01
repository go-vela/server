// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// MaxBuildLimit is a middleware function that attaches the defaultLimit
// to enable the server to override the max build limit.
func MaxBuildLimit(maxBuildLimit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("maxBuildLimit", maxBuildLimit)
		c.Next()
	}
}
