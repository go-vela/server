// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/perm"
)

// LogServiceHandlers is a function that extends the provided base router group
// with the API handlers for service logs functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
func LogServiceHandlers(base *gin.RouterGroup) {
	// Logs endpoints
	logs := base.Group("/logs")
	{
		logs.POST("", perm.MustAdmin(), api.CreateServiceLog)
		logs.GET("", perm.MustRead(), api.GetServiceLog)
		logs.PUT("", perm.MustWrite(), api.UpdateServiceLog)
		logs.DELETE("", perm.MustPlatformAdmin(), api.DeleteServiceLog)
	} // end of logs endpoints
}

// LogStepHandlers is a function that extends the provided base router group
// with the API handlers for step logs functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
func LogStepHandlers(base *gin.RouterGroup) {
	// Logs endpoints
	logs := base.Group("/logs")
	{
		logs.POST("", perm.MustAdmin(), api.CreateStepLog)
		logs.GET("", perm.MustRead(), api.GetStepLog)
		logs.PUT("", perm.MustWrite(), api.UpdateStepLog)
		logs.DELETE("", perm.MustPlatformAdmin(), api.DeleteStepLog)
	} // end of logs endpoints
}
