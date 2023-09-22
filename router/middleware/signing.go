// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// QueueSigningPrivateKey is a middleware function that attaches the private key used
// to sign items that are pushed to the queue.
func QueueSigningPrivateKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("queue.private-key", key)
		c.Next()
	}
}
