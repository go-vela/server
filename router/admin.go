// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/go-vela/server/api/admin"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/gin-gonic/gin"
)

// AdminHandlers is a function that extends the provided base router group
// with the API handlers for admin functionality.
//
// GET    /api/v1/admin/builds
// GET    /api/v1/admin/repos
// GET    /api/v1/admin/secrets
// GET    /api/v1/admin/steps
// GET    /api/v1/admin/users
func AdminHandlers(base *gin.RouterGroup) {

	// Admin endpoints
	_admin := base.Group("/admin", perm.MustPlatformAdmin())
	{
		_admin.GET("/builds", admin.AllBuilds)
		_admin.GET("/repos", admin.AllRepos)
		_admin.GET("/secrets", admin.AllSecrets)
		_admin.GET("/steps", admin.AllSteps)
		_admin.GET("/users", admin.AllUsers)
	} // end of admin endpoints

}
