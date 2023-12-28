// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with step
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/service"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/perm"
	smiddleware "github.com/go-vela/server/router/middleware/service"
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
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service/logs .
func ServiceHandlers(base *gin.RouterGroup) {
	// Services endpoints
	services := base.Group("/services")
	{
		services.POST("", perm.MustPlatformAdmin(), middleware.Payload(), service.CreateService)
		services.GET("", perm.MustRead(), service.ListServices)

		// Service endpoints
		s := services.Group("/:service", smiddleware.Establish())
		{
			s.GET("", perm.MustRead(), service.GetService)
			s.PUT("", perm.MustBuildAccess(), middleware.Payload(), service.UpdateService)
			s.DELETE("", perm.MustPlatformAdmin(), service.DeleteService)

			// Log endpoints
			LogServiceHandlers(s)
		} // end of service endpoints
	} // end of services endpoints
}
