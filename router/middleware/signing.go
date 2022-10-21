// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// QueueSigningPrivateKey is a middleware function that attaches the private key used
// to sign items that are pushed to the queue.
func QueueSigningPrivateKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("queue.signing.private-key", key)
		c.Next()
	}
}
