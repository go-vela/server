// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/queue"
)

// Queue is a middleware function that initializes the queue and
// attaches to the context of every http.Request.
func Queue(q queue.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		queue.ToContext(c, q)
		c.Next()
	}
}
