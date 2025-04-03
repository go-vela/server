// SPDX-License-Identifier: Apache-2.0

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
// to open signed items that are pushed to the queue.
func QueueSigningPublicKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("public-key", key)
		c.Next()
	}
}

// QueueAddress is a middleware function that attaches the queue address used
// to open the connection to the queue.
func QueueAddress(address string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("queue-address", address)
		c.Next()
	}
}

// StorageAccessKey is a middleware function that attaches the access key used
// to open the connection to the storage.
func StorageAccessKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("access-key", key)
		c.Next()
	}
}

// StorageSecretKey is a middleware function that attaches the secret key used
// to open the connection to the storage.
func StorageSecretKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("secret-key", key)
		c.Next()
	}
}

// StorageAddress is a middleware function that attaches the storage address used
// to open the connection to the storage.
func StorageAddress(address string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage-address", address)
		c.Next()
	}
}

// StorageBucket is a middleware function that attaches the bucket name used
// to open the connection to the storage.
func StorageBucket(bucket string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage-bucket", bucket)
		c.Next()
	}
}

// StorageEnable is a middleware function that sets a flag in the context
// to determined if storage is enabled.
func StorageEnable(enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage-enable", enabled)
		c.Next()
	}
}
