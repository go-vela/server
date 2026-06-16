// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// MaxOrgBuildLimit is a middleware function that attaches the maxOrgBuildLimit
// to enable the server to override the max org build limit.
func MaxOrgBuildLimit(maxOrgBuildLimit int32) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("maxOrgBuildLimit", maxOrgBuildLimit)
		c.Next()
	}
}
