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
		c.Set("queue.private-key", key)
		c.Next()
	}
}

// QueueSigningPublicKey is a middleware function that attaches the public key used
// to opened signed items that are pushed to the queue.
func QueueSigningPublicKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("public-key", key)
		c.Next()
	}
}

// QueueAddress is a middleware function that attaches the queue address used
// to open the connection to the queue.
func QueueAddress(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("queue-address", key)
		c.Next()
	}
}
