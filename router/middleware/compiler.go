// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/compiler/compiler"
)

// Compiler is a middleware function that initializes the compiler and
// attaches to the context of every http.Request.
func Compiler(cli compiler.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		compiler.ToContext(c, cli)
		c.Next()
	}
}
