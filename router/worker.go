// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/worker"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/perm"
	wmiddleware "github.com/go-vela/server/router/middleware/worker"
)

// WorkerHandlers is a function that extends the provided base router group
// with the API handlers for worker functionality.
//
// POST   /api/v1/users
// GET    /api/v1/workers
// GET    /api/v1/workers/:worker
// PUT    /api/v1/workers/:worker
// POST   /api/v1/workers/:worker/refresh
// DELETE /api/v1/workers/:worker .
func WorkerHandlers(base *gin.RouterGroup) {
	// Workers endpoints
	_workers := base.Group("/workers")
	{
		_workers.POST("", perm.MustWorkerRegisterToken(), middleware.Payload(), worker.CreateWorker)
		_workers.GET("", worker.ListWorkers)

		// Worker endpoints
		_worker := _workers.Group("/:worker")
		{
			_worker.GET("", wmiddleware.Establish(), worker.GetWorker)
			_worker.PUT("", perm.MustPlatformAdmin(), perm.MustWorkerAuthToken(), wmiddleware.Establish(), worker.UpdateWorker)
			_worker.POST("/refresh", perm.MustWorkerAuthToken(), wmiddleware.Establish(), worker.Refresh)
			_worker.DELETE("", perm.MustPlatformAdmin(), wmiddleware.Establish(), worker.DeleteWorker)
		} // end of worker endpoints
	} // end of workers endpoints
}
