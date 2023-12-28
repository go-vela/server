// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
)

// Database is a middleware function that initializes the database and
// attaches to the context of every http.Request.
func Database(d database.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		database.ToContext(c, d)
		c.Next()
	}
}
