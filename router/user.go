// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/build"
	"github.com/go-vela/server/api/user"
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
// GET    /api/v1/user
// PUT    /api/v1/user
// GET    /api/v1/user/source/repos
// GET    /api/v1/user/builds
// POST   /api/v1/user/token
// DELETE /api/v1/user/token .
func UserHandlers(base *gin.RouterGroup) {
	// Users endpoints
	_users := base.Group("/users")
	{
		_users.POST("", perm.MustPlatformAdmin(), user.CreateUser)
		_users.GET("", user.ListUsers)
		_users.GET("/:user", perm.MustPlatformAdmin(), user.GetUser)
		_users.PUT("/:user", perm.MustPlatformAdmin(), user.UpdateUser)
		_users.DELETE("/:user", perm.MustPlatformAdmin(), user.DeleteUser)
	} // end of users endpoints

	// User endpoints
	_user := base.Group("/user")
	{
		_user.GET("", user.GetCurrentUser)
		_user.PUT("", user.UpdateCurrentUser)
		_user.GET("/source/repos", user.GetSourceRepos)
		_user.GET("/builds", build.ListBuildsForSender)
		_user.POST("/token", user.CreateToken)
		_user.DELETE("/token", user.DeleteToken)
	} // end of user endpoints
}
