// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/storage"
)

// StorageHandlers is a function that extends the provided base router group
// with the API handlers for storage functionality.
//
// GET   /api/v1/storage/info .
func StorageHandlers(base *gin.RouterGroup) {
	// Storage endpoints
	_storage := base.Group("/storage")
	{
		_storage.GET("/:bucket/objects", storage.ListObjects)
		_storage.GET("/:bucket/names", storage.ListObjectNames)
		_storage.GET("/:bucket/:org/:repo/builds/:build/names", storage.ListBuildObjectNames)
		_storage.GET("/sts/:org/:repo/:build", storage.GetSTSCreds)
	} // end of storage endpoints
}
