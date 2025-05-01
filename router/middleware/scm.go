// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/router/middleware/settings"
	"github.com/go-vela/server/scm"
)

// Scm is a middleware function that initializes the scm and
// attaches to the context of every http.Request.
func Scm(scmService scm.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := settings.FromContext(c)
		scmService.SetSettings(s)

		scm.WithGinContext(c, scmService)

		c.Next()
	}
}
