// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/perm"
)

// UserHandlers is a function that extends the provided base router group
// with the API handlers for user functionality.
//
// POST   /api/v1/users
// GET    /api/v1/users
// GET    /api/v1/users/:user
// PUT    /api/v1/users/:user
// DELETE /api/v1/users/:user
// GET    /api/v1/users/:user/source/repos
// POST   /api/v1/users/:user/token
// DELETE /api/v1/users/:user/token
// GET /api/v1/user
// PUT /api/v1/user
// GET /api/v1/source/repos
// POST /api/v1/token
// DELETE /api/v1/token
func UserHandlers(base *gin.RouterGroup) {
	// Users endpoints
	users := base.Group("/users")
	{
		users.POST("", perm.MustPlatformAdmin(), api.CreateUser)
		users.GET("", api.GetUsers)
		users.GET("/:user", perm.MustPlatformAdmin(), api.GetUser)
		users.PUT("/:user", perm.MustPlatformAdmin(), api.UpdateUser)
		users.DELETE("/:user", perm.MustPlatformAdmin(), api.DeleteUser)
	} // end of users endpoints

	// User endpoints
	user := base.Group("/user")
	{
		user.GET("", api.GetCurrentUser)
		user.PUT("", api.UpdateCurrentUser)
		user.GET("/source/repos", api.GetUserSourceRepos)
		user.POST("/token", api.CreateToken)
		user.DELETE("/token", api.DeleteToken)
	} // end of user endpoints
}
