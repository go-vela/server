// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/queue"
)

// Queue is a middleware function that initializes the queue and
// attaches to the context of every http.Request.
func Queue(q queue.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		queue.WithGinContext(c, q)
		c.Next()
	}
}
