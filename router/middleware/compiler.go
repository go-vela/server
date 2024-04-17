// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/router/middleware/settings"
)

// Compiler is a middleware function that initializes the compiler and
// attaches to the context of every http.Request.
func Compiler(comp compiler.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := settings.FromContext(c)
		comp.SetSettings(s)

		compiler.WithGinContext(c, comp)

		c.Next()
	}
}
