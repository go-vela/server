// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/org"

	"github.com/go-vela/server/api/deployment"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/repo"
)

// DeploymentHandlers is a function that extends the provided base router group
// with the API handlers for deployment functionality.
//
// POST   /api/v1/deployments/:org/:repo
// GET    /api/v1/deployments/:org/:repo
// GET    /api/v1/deployments/:org/:repo/:deployment .
func DeploymentHandlers(base *gin.RouterGroup) {
	// Deployments endpoints
	deployments := base.Group("/deployments/:org/:repo", org.Establish(), repo.Establish())
	{
		deployments.POST("", perm.MustWrite(), deployment.CreateDeployment)
		deployments.GET("", perm.MustRead(), deployment.ListDeployments)
		deployments.GET("/:deployment", perm.MustRead(), deployment.GetDeployment)
	} // end of deployments endpoints
}
