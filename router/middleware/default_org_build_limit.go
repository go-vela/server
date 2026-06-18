// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// DefaultOrgBuildLimit is a middleware function that attaches the defaultOrgBuildLimit
// to enable the server to override the default org build limit.
func DefaultOrgBuildLimit(defaultOrgBuildLimit int32) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("defaultOrgBuildLimit", defaultOrgBuildLimit)
		c.Next()
	}
}
