// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/hook"
	"github.com/go-vela/server/router/middleware/org"
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
// POST   /api/v1/hooks/:org/:repo/:hook/redeliver .
func HookHandlers(base *gin.RouterGroup) {
	// Hooks endpoints
	_hooks := base.Group("/hooks/:org/:repo", org.Establish(), repo.Establish())
	{
		_hooks.POST("", perm.MustPlatformAdmin(), hook.CreateHook)
		_hooks.GET("", perm.MustRead(), hook.ListHooks)
		_hooks.GET("/:hook", perm.MustRead(), hook.GetHook)
		_hooks.PUT("/:hook", perm.MustPlatformAdmin(), hook.UpdateHook)
		_hooks.DELETE("/:hook", perm.MustPlatformAdmin(), hook.DeleteHook)
		_hooks.POST("/:hook/redeliver", perm.MustWrite(), hook.RedeliverHook)
	} // end of hooks endpoints
}
