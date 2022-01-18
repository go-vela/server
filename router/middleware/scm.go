// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/scm"
)

// Scm is a middleware function that initializes the scm and
// attaches to the context of every http.Request.
func Scm(s scm.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		scm.ToContext(c, s)
		c.Next()
	}
}
