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
// GET   /api/v1/storage/info .
func StorageHandlers(base *gin.RouterGroup) {
	// Storage endpoints
	_storage := base.Group("/storage")
	{
		_storage.GET("/info", perm.MustWorkerRegisterToken(), storage.Info)
		_storage.GET("/:bucket/objects", storage.ListObjects)
	} // end of storage endpoints
}
