// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// DefaultLimit is a middleware function that attaches the defaultLimit
// to enable the server to override the default build limit.
func DefaultLimit(defaultLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("defaultLimit", defaultLimit)
		c.Next()
	}
}
