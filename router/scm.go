// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/scm"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
)

// ScmHandlers is a function that extends the provided base router group
// with the API handlers for source code management functionality.
//
// PATCH   /api/v1/scm/orgs/:org/sync
// PATCH   /api/v1/scm/repos/:org/:repo/sync .
func ScmHandlers(base *gin.RouterGroup) {
	// SCM orgs endpoints
	orgs := base.Group("/scm/orgs")
	{
		// SCM org endpoints
		org := orgs.Group("/:org", org.Establish())
		{
			org.PATCH("/sync", scm.SyncReposForOrg)
		} // end of SCM org endpoints
	} // end of SCM orgs endpoints

	// SCM repos endpoints
	repos := base.Group("/scm/repos")
	{
		// SCM repo endpoints
		repo := repos.Group("/:org/:repo", org.Establish(), repo.Establish())
		{
			repo.PATCH("/sync", scm.SyncRepo)
		} // end of SCM repo endpoints
	} // end of SCM repos endpoints
}
