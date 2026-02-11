// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/storage"
	"github.com/go-vela/server/router/middleware/perm"
)

// StorageHandlers is a function that extends the provided base router group
// with the API handlers for storage functionality.
//
// GET   /api/v1/repos/:org/:repo/builds/:build/storage/sts
// GET   /api/v1/repos/:org/:repo/builds/:build/storage/:bucket/names.
func StorageHandlers(base *gin.RouterGroup) {
	// Storage endpoints
	_storage := base.Group("/storage")
	{
		_storage.GET("/:bucket/names", perm.MustRead(), storage.ListBuildObjectNames)
		_storage.GET("/sts", perm.MustBuildAccess(), storage.GetSTSCreds)
	} // end of storage endpoints
}
