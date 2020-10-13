// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/worker"
)

// WorkerHandlers is a function that extends the provided base router group
// with the API handlers for worker functionality.
//
// POST   /api/v1/users
func WorkerHandlers(base *gin.RouterGroup) {
	// Workers endpoints
	workers := base.Group("/workers")
	{
		workers.POST("", perm.MustPlatformAdmin(), middleware.Payload(), api.CreateWorker)
		workers.GET("", api.GetWorkers)

		// Worker endpoints
		w := workers.Group("/:worker")
		{
			w.GET("", worker.Establish(), api.GetWorker)
			w.PUT("", perm.MustPlatformAdmin(), worker.Establish(), api.UpdateWorker)
			w.DELETE("", perm.MustPlatformAdmin(), worker.Establish(), api.DeleteWorker)
		} // end of worker endpoints
	} // end of workers endpoints
}
