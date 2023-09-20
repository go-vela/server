// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/admin"
	"github.com/go-vela/server/router/middleware/perm"
)

// AdminHandlers is a function that extends the provided base router group
// with the API handlers for admin functionality.
//
// GET    /api/v1/admin/builds/queue
// GET    /api/v1/admin/build/:id
// PUT    /api/v1/admin/build
// PUT    /api/v1/admin/clean
// PUT    /api/v1/admin/deployment
// PUT    /api/v1/admin/hook
// PUT    /api/v1/admin/repo
// PUT    /api/v1/admin/secret
// PUT    /api/v1/admin/service
// PUT    /api/v1/admin/step
// PUT    /api/v1/admin/user
// POST   /api/v1/admin/workers/:worker/register.
func AdminHandlers(base *gin.RouterGroup) {
	// Admin endpoints
	_admin := base.Group("/admin", perm.MustPlatformAdmin())
	{
		// Admin build queue endpoint
		_admin.GET("/builds/queue", admin.AllBuildsQueue)

		// Admin build endpoint
		_admin.PUT("/build", admin.UpdateBuild)

		// Admin clean endpoint
		_admin.PUT("/clean", admin.CleanResources)

		// Admin deployment endpoint
		_admin.PUT("/deployment", admin.UpdateDeployment)

		// Admin hook endpoint
		_admin.PUT("/hook", admin.UpdateHook)

		// Admin repo endpoint
		_admin.PUT("/repo", admin.UpdateRepo)

		// Admin secret endpoint
		_admin.PUT("/secret", admin.UpdateSecret)

		// Admin service endpoint
		_admin.PUT("/service", admin.UpdateService)

		// Admin step endpoint
		_admin.PUT("/step", admin.UpdateStep)

		// Admin user endpoint
		_admin.PUT("/user", admin.UpdateUser)

		// Admin worker endpoint
		_admin.POST("/workers/:worker/register", admin.RegisterToken)
	} // end of admin endpoints
}
