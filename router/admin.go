// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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
// GET    /api/v1/admin/builds
// PUT    /api/v1/admin/build
// GET    /api/v1/admin/hooks
// PUT    /api/v1/admin/hook
// GET    /api/v1/admin/repos
// PUT    /api/v1/admin/repo
// GET    /api/v1/admin/secrets
// PUT    /api/v1/admin/secret
// GET    /api/v1/admin/services
// PUT    /api/v1/admin/service
// GET    /api/v1/admin/steps
// PUT    /api/v1/admin/step
// GET    /api/v1/admin/users
// PUT    /api/v1/admin/user
func AdminHandlers(base *gin.RouterGroup) {
	// Admin endpoints
	_admin := base.Group("/admin", perm.MustPlatformAdmin())
	{
		// Admin build endpoints
		_admin.GET("/builds", admin.AllBuilds)
		_admin.PUT("/build", admin.UpdateBuild)

		// Admin hook endpoints
		_admin.GET("/hooks", admin.AllHooks)
		_admin.PUT("/hook", admin.UpdateHook)

		// Admin repo endpoints
		_admin.GET("/repos", admin.AllRepos)
		_admin.PUT("/repo", admin.UpdateRepo)

		// Admin secret endpoints
		_admin.GET("/secrets", admin.AllSecrets)
		_admin.PUT("/secret", admin.UpdateSecret)

		// Admin service endpoints
		_admin.GET("/services", admin.AllServices)
		_admin.PUT("/service", admin.UpdateService)

		// Admin step endpoints
		_admin.GET("/steps", admin.AllSteps)
		_admin.PUT("/step", admin.UpdateStep)

		// Admin user endpoints
		_admin.GET("/users", admin.AllUsers)
		_admin.PUT("/user", admin.UpdateUser)
	} // end of admin endpoints
}
