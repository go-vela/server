// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/log"
	"github.com/go-vela/server/router/middleware/perm"
)

// LogServiceHandlers is a function that extends the provided base router group
// with the API handlers for service logs functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service/logs .
func LogServiceHandlers(base *gin.RouterGroup) {
	// Logs endpoints
	logs := base.Group("/logs")
	{
		logs.POST("", perm.MustAdmin(), log.CreateServiceLog)
		logs.GET("", perm.MustRead(), log.GetServiceLog)
		logs.PUT("", perm.MustBuildAccess(), log.UpdateServiceLog)
		logs.DELETE("", perm.MustPlatformAdmin(), log.DeleteServiceLog)
	} // end of logs endpoints
}

// LogStepHandlers is a function that extends the provided base router group
// with the API handlers for step logs functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs .
func LogStepHandlers(base *gin.RouterGroup) {
	// Logs endpoints
	logs := base.Group("/logs")
	{
		logs.POST("", perm.MustAdmin(), log.CreateStepLog)
		logs.GET("", perm.MustRead(), log.GetStepLog)
		logs.PUT("", perm.MustBuildAccess(), log.UpdateStepLog)
		logs.DELETE("", perm.MustPlatformAdmin(), log.DeleteStepLog)
	} // end of logs endpoints
}
