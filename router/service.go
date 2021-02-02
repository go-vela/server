// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// nolint: dupl // ignore similar code with step
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/service"
)

// ServiceHandlers is a function that extends the provided base router group
// with the API handlers for service functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/services
// GET    /api/v1/repos/:org/:repo/builds/:build/services
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service
// POST   /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
func ServiceHandlers(base *gin.RouterGroup) {
	// Services endpoints
	services := base.Group("/services")
	{
		services.POST("", perm.MustPlatformAdmin(), middleware.Payload(), api.CreateService)
		services.GET("", perm.MustRead(), api.GetServices)

		// Service endpoints
		service := services.Group("/:service", service.Establish())
		{
			service.GET("", perm.MustRead(), api.GetService)
			service.PUT("", perm.MustPlatformAdmin(), middleware.Payload(), api.UpdateService)
			service.DELETE("", perm.MustPlatformAdmin(), api.DeleteService)

			// Log endpoints
			LogServiceHandlers(service)
		} // end of service endpoints
	} // end of services endpoints
}
