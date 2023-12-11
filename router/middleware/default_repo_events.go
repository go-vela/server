// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// DefaultRepoEvents is a middleware function that attaches the defaultRepoEvents
// to enable the server to override the default repo event.
func DefaultRepoEvents(defaultRepoEvents []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("defaultRepoEvents", defaultRepoEvents)
		c.Next()
	}
}

// DefaultRepoEventsMask is a middleware function that attaches the defaultRepoEventsMask
// to enable the server to override the default repo events using a mask.
func DefaultRepoEventsMask(defaultRepoEventsMask int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("defaultRepoEventsMask", defaultRepoEventsMask)
		c.Next()
	}
}
