// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/repo"
)

// HookHandlers is a function that extends the provided base router group
// with the API handlers for hook functionality.
//
// POST   /api/v1/hooks/:org/:repo
// GET    /api/v1/hooks/:org/:repo
// GET    /api/v1/hooks/:org/:repo/:hook
// PUT    /api/v1/hooks/:org/:repo/:hook
// DELETE /api/v1/hooks/:org/:repo/:hook
func HookHandlers(base *gin.RouterGroup) {
	// Hooks endpoints
	hooks := base.Group("/hooks/:org/:repo", repo.Establish())
	{
		hooks.POST("", perm.MustPlatformAdmin(), api.CreateHook)
		hooks.GET("", perm.MustRead(), api.GetHooks)
		hooks.GET("/:hook", perm.MustRead(), api.GetHook)
		hooks.PUT("/:hook", perm.MustPlatformAdmin(), api.UpdateHook)
		hooks.DELETE("/:hook", perm.MustPlatformAdmin(), api.DeleteHook)
	} // end of hooks endpoints
}
