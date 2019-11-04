// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/gin-gonic/gin"
)

// BuildHandlers is a function that extends the provided base router group
// with the API handlers for build functionality.
//
// POST   /api/v1/repos/:org/:repo/builds
// GET    /api/v1/repos/:org/:repo/builds
// POST   /api/v1/repos/:org/:repo/builds/:build
// GET    /api/v1/repos/:org/:repo/builds/:build
// PUT    /api/v1/repos/:org/:repo/builds/:build
// DELETE /api/v1/repos/:org/:repo/builds/:build
// GET    /api/v1/repos/:org/:repo/builds/:build/logs
// POST   /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step
// POST   /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
func BuildHandlers(base *gin.RouterGroup) {

	// Builds endpoints
	builds := base.Group("/builds")
	{
		builds.POST("", perm.MustAdmin(), middleware.Payload(), api.CreateBuild)
		builds.GET("", perm.MustRead(), api.GetBuilds)

		// Build endpoints
		build := builds.Group("/:build", build.Establish())
		{
			build.POST("", perm.MustWrite(), api.RestartBuild)
			build.GET("", perm.MustRead(), api.GetBuild)
			build.PUT("", perm.MustWrite(), middleware.Payload(), api.UpdateBuild)
			build.DELETE("", perm.MustPlatformAdmin(), api.DeleteBuild)
			build.GET("/logs", perm.MustRead(), api.GetBuildLogs)

			// Step endpoints
			// * Log endpoints
			StepHandlers(build)
			ServiceHandlers(build)

		} // end of build endpoints

	} // end of builds endpoints

}
