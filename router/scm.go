// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
)

// ScmHandlers is a function that extends the provided base router group
// with the API handlers for source code management functionality.
//
// GET    /api/v1/scm/orgs/:org/sync
// GET    /api/v1/scm/repos/:org/:repo/sync .
func ScmHandlers(base *gin.RouterGroup) {
	repos := base.Group("/scm/repos")
	orgs := base.Group("/scm/orgs")
	{
		org := orgs.Group("/:org", org.Establish())
		fmt.Printf("ORG: %s", org.BasePath())
		{
			org.GET("/sync", api.SyncRepos)
		} // end of org endpoints
		// Repo endpoints
		repo := repos.Group("/:org/:repo", repo.Establish())
		{
			repo.GET("/sync", api.SyncRepo)
		} // end of repo endpoints
	} // end of scm endpoints
}
