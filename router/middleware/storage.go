// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/storage"
)

// Storage is a middleware function that initializes the object storage and
// attaches to the context of every http.Request.
func Storage(q storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// attach the object storage to the context
		storage.WithGinContext(c, q)

		c.Next()
	}
}
