// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// MaxBuildLimit is a middleware function that attaches the defaultLimit
// to enable the server to override the max build limit.
func MaxBuildLimit(maxBuildLimit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("maxBuildLimit", maxBuildLimit)
		c.Next()
	}
}
