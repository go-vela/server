// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/storage"
)

// Storage is a middleware function that initializes the object storage and
// attaches to the context of every http.Request.
func Storage(q storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// attach the object storage to the context
		storage.WithGinContext(c, q)

		c.Next()
	}
}

// StorageAccessKey is a middleware function that attaches the access key used
// to open the connection to the storage.
func StorageAccessKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage-access-key", key)
		c.Next()
	}
}

// StorageSecretKey is a middleware function that attaches the secret key used
// to open the connection to the storage.
func StorageSecretKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage-secret-key", key)
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
