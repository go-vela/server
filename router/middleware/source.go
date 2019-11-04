// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/go-vela/server/source"
	"github.com/gin-gonic/gin"
)

// Source is a middleware function that initializes the source and
// attaches to the context of every http.Request.
func Source(s source.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		source.ToContext(c, s)
		c.Next()
	}
}
