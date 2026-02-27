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
// GET   /api/v1/repos/:org/:repo/builds/:build/storage/:name/presigned-put
// GET   /api/v1/repos/:org/:repo/builds/:build/storage/.
func StorageHandlers(base *gin.RouterGroup) {
	// Storage endpoints
	_storage := base.Group("/storage")
	{
		_storage.GET("/", perm.MustRead(), storage.ListBuildObjectNames)
		_storage.GET("/:name/presigned-put", perm.MustBuildAccess(), storage.GetPresignedPutURL)
	} // end of storage endpoints
}
