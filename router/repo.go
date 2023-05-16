// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/api/repo"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/perm"
	rmiddleware "github.com/go-vela/server/router/middleware/repo"
)

// RepoHandlers is a function that extends the provided base router group
// with the API handlers for repo functionality.
//
// POST   /api/v1/repos
// GET    /api/v1/repos
// GET    /api/v1/repos/:org
// GET    /api/v1/repos/:org/builds
// GET    /api/v1/repos/:org/:repo
// PUT    /api/v1/repos/:org/:repo
// DELETE /api/v1/repos/:org/:repo
// PATCH  /api/v1/repos/:org/:repo/repair
// PATCH  /api/v1/repos/:org/:repo/chown
// POST   /api/v1/repos/:org/:repo/builds
// GET    /api/v1/repos/:org/:repo/builds
// POST   /api/v1/repos/:org/:repo/builds/:build
// GET    /api/v1/repos/:org/:repo/builds/:build
// PUT    /api/v1/repos/:org/:repo/builds/:build
// DELETE /api/v1/repos/:org/:repo/builds/:build
// DELETE /api/v1/repos/:org/:repo/builds/:build/cancel
// GET    /api/v1/repos/:org/:repo/builds/:build/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/token
// POST   /api/v1/repos/:org/:repo/builds/:build/services
// GET    /api/v1/repos/:org/:repo/builds/:build/services
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service
// POST   /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/services/:service/logs
// POST   /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step
// POST   /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs .
func RepoHandlers(base *gin.RouterGroup) {
	// Repos endpoints
	_repos := base.Group("/repos")
	{
		_repos.POST("", middleware.Payload(), repo.CreateRepo)
		_repos.GET("", repo.ListRepos)

		// Org endpoints
		org := _repos.Group("/:org", org.Establish())
		{
			org.GET("", repo.ListReposForOrg)
			org.GET("/builds", api.GetOrgBuilds)

			// Repo endpoints
			_repo := org.Group("/:repo", rmiddleware.Establish())
			{
				_repo.GET("", perm.MustRead(), repo.GetRepo)
				_repo.PUT("", perm.MustAdmin(), middleware.Payload(), repo.UpdateRepo)
				_repo.DELETE("", perm.MustAdmin(), repo.DeleteRepo)
				_repo.PATCH("/repair", perm.MustAdmin(), repo.RepairRepo)
				_repo.PATCH("/chown", perm.MustAdmin(), repo.ChownRepo)

				// Build endpoints
				// * Service endpoints
				//   * Log endpoints
				// * Step endpoints
				//   * Log endpoints
				BuildHandlers(_repo)
			} // end of repo endpoints
		} // end of org endpoints
	} // end of repos endpoints
}
