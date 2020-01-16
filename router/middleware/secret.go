// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/secret"
)

// Secret is a middleware function that attaches the secret used for
// server <-> agent communication to the context of every http.Request.
func Secret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("secret", secret)
		c.Next()
	}
}

// Secrets is a middleware function that initializes the secret engines and
// attaches to the context of every http.Request.
func Secrets(secrets map[string]secret.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range secrets {
			secret.ToContext(c, k, v)
		}
		c.Next()
	}
}
