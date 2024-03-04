// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/build"
	"github.com/go-vela/server/api/log"
	"github.com/go-vela/server/router/middleware"
	bmiddleware "github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/executors"
	"github.com/go-vela/server/router/middleware/perm"
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
// POST   /api/v1/repos/:org/:repo/builds/:build/approve
// DELETE /api/v1/repos/:org/:repo/builds/:build/cancel
// GET    /api/v1/repos/:org/:repo/builds/:build/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/token
// GET    /api/v1/repos/:org/:repo/builds/:build/executable
// POST   /api/v1/repos/:org/:repo/builds/:build/services
// GET    /api/v1/repos/:org/:repo/builds/:build/services
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service
// POST   /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// POST   /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step
// POST   /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs .
func BuildHandlers(base *gin.RouterGroup) {
	// Builds endpoints
	builds := base.Group("/builds")
	{
		builds.POST("", perm.MustAdmin(), middleware.Payload(), build.CreateBuild)
		builds.GET("", perm.MustRead(), build.ListBuildsForRepo)

		// Build endpoints
		b := builds.Group("/:build", bmiddleware.Establish())
		{
			b.POST("", perm.MustWrite(), build.RestartBuild)
			b.GET("", perm.MustRead(), build.GetBuild)
			b.PUT("", perm.MustBuildAccess(), middleware.Payload(), build.UpdateBuild)
			b.DELETE("", perm.MustPlatformAdmin(), build.DeleteBuild)
			b.POST("/approve", perm.MustAdmin(), build.ApproveBuild)
			b.DELETE("/cancel", executors.Establish(), perm.MustWrite(), build.CancelBuild)
			b.GET("/logs", perm.MustRead(), log.ListLogsForBuild)
			b.GET("/token", perm.MustWorkerAuthToken(), build.GetBuildToken)
			b.GET("/graph", perm.MustRead(), build.GetBuildGraph)
			b.GET("/executable", perm.MustBuildAccess(), build.GetBuildExecutable)

			// Service endpoints
			// * Log endpoints
			ServiceHandlers(b)

			// Step endpoints
			// * Log endpoints
			StepHandlers(b)
		} // end of build endpoints
	} // end of builds endpoints
}
