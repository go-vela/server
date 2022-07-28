// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
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
	// SCM orgs endpoints
	orgs := base.Group("/scm/orgs")
	{
		// SCM org endpoints
		org := orgs.Group("/:org", org.Establish())
		{
			org.GET("/sync", api.SyncRepos)
		} // end of SCM org endpoints
	} // end of SCM orgs endpoints

	// SCM repos endpoints
	repos := base.Group("/scm/repos")
	{
		// SCM repo endpoints
		repo := repos.Group("/:org/:repo", org.Establish(), repo.Establish())
		{
			repo.GET("/sync", api.SyncRepo)
		} // end of SCM repo endpoints
	} // end of SCM repos endpoints

	// SCM issue endpoints
	issue := base.Group("/scm/issue")
	{
		issue.POST("/create", api.CreateIssue)
	}
}
