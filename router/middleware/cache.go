// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/cache"
)

// Cache is a middleware function that initializes the cache and
// attaches to the context of every http.Request.
func Cache(cs cache.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cache", cs)

		c.Next()
	}
}
