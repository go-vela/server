// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/types/settings"
	sMiddleware "github.com/go-vela/server/router/middleware/settings"
)

// Settings is a middleware function that attaches settings
// to the context of every http.Request.
func Settings(s *settings.Platform) gin.HandlerFunc {
	return func(c *gin.Context) {
		sMiddleware.ToContext(c, s)

		c.Next()
	}
}
