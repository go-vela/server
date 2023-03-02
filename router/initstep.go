// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with service
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	initStepApi "github.com/go-vela/server/api/initstep"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/initstep"
	"github.com/go-vela/server/router/middleware/perm"
)

// InitStepHandlers is a function that extends the provided base router group
// with the API handlers for initstep functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/initsteps
// GET    /api/v1/repos/:org/:repo/builds/:build/initsteps
// GET    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep
// PUT    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep
// DELETE /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep
// POST   /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
// POST   /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/stream .
func InitStepHandlers(base *gin.RouterGroup) {
	// InitSteps endpoints
	initsteps := base.Group("/initsteps")
	{
		initsteps.POST("", perm.MustPlatformAdmin(), middleware.Payload(), initStepApi.CreateInitStep)
		initsteps.GET("", perm.MustRead(), initStepApi.ListInitSteps)

		// InitStep endpoints
		initstep := initsteps.Group("/:initstep", initstep.Establish())
		{
			initstep.GET("", perm.MustRead(), initStepApi.GetInitStep)
			initstep.PUT("", perm.MustBuildAccess(), middleware.Payload(), initStepApi.UpdateInitStep)
			initstep.DELETE("", perm.MustPlatformAdmin(), initStepApi.DeleteInitStep)

			initstep.POST("/stream", perm.MustPlatformAdmin(), api.PostInitStepStream)

			// Log endpoints
			LogInitStepHandlers(initstep)
		} // end of initstep endpoints
	} // end of initsteps endpoints
}
