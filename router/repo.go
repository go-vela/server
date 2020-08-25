// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/repo"
)

// RepoHandlers is a function that extends the provided base router group
// with the API handlers for repo functionality.
//
// POST   /api/v1/repos
// GET    /api/v1/repos
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
// GET    /api/v1/repos/:org/:repo/builds/:build/logs
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
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
func RepoHandlers(base *gin.RouterGroup) {
	// Repos endpoints
	repos := base.Group("/repos")
	{
		repos.POST("", middleware.Payload(), api.CreateRepo)
		repos.GET("", api.GetRepos)

		// Repo endpoints
		repo := repos.Group("/:org/:repo", repo.Establish()) //[here] instantiate(spellcheck) router gin group "repo"
		{                                                    //[here] Adds things to router group
			repo.GET("", perm.MustRead(), api.GetRepo)
			repo.PUT("", perm.MustAdmin(), middleware.Payload(), api.UpdateRepo)
			repo.DELETE("", perm.MustAdmin(), api.DeleteRepo)
			repo.PATCH("/repair", perm.MustAdmin(), api.RepairRepo)
			repo.PATCH("/chown", perm.MustAdmin(), api.ChownRepo)

			// Build endpoints
			// * Service endpoints
			//   * Log endpoints
			// * Step endpoints
			//   * Log endpoints
			BuildHandlers(repo) //[here] step 2
		} // end of repo endpoints
	} // end of repos endpoints
}
