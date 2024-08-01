// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecureCookie determines whether or not incoming webhooks are validated coming from Github
// This is primarily intended for local development.
func SecureCookie(secure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("securecookie", secure)
		c.Next()
	}
}
