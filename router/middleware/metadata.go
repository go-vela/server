// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/types"
)

// Metadata is a middleware function that attaches the metadata
// to the context of every http.Request.
func Metadata(m *types.Metadata) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("metadata", m)
		c.Next()
	}
}
