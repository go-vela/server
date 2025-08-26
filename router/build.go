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

// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs.
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
			b.GET("/id_token", perm.MustIDRequestToken(), build.GetIDToken)
			b.GET("/id_request_token", perm.MustBuildAccess(), build.GetIDRequestToken)
			b.GET("/graph", perm.MustRead(), build.GetBuildGraph)
			b.GET("/executable", perm.MustBuildAccess(), build.GetBuildExecutable)

			// Service endpoints
			// * Log endpoints
			ServiceHandlers(b)

			// Step endpoints
			// * Log endpoints
			StepHandlers(b)

			// Test attachment endpoints
			TestAttachmentHandlers(b)

			// Test report endpoints
			TestReportHandlers(b)
		} // end of build endpoints
	} // end of builds endpoints
}
