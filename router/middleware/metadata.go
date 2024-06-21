// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/internal"
)

// Metadata is a middleware function that attaches the metadata
// to the context of every http.Request.
func Metadata(m *internal.Metadata) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("metadata", m)
		c.Next()
	}
}
