// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Worker is a middleware function that attaches the worker interval
// to determine which workers are active.
func Worker(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("worker_active_interval", duration)
		c.Next()
	}
}
