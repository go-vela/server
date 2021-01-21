// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/repo"
)

// DeploymentHandlers is a function that extends the provided base router group
// with the API handlers for deployment functionality.
//
// POST   /api/v1/deployments/:org/:repo
// GET    /api/v1/deployments/:org/:repo
// GET    /api/v1/deployments/:org/:repo/:deployment
func DeploymentHandlers(base *gin.RouterGroup) {
	// Deployments endpoints
	deployments := base.Group("/deployments/:org/:repo", repo.Establish())
	{
		deployments.POST("", perm.MustWrite(), api.CreateDeployment)
		deployments.GET("", perm.MustRead(), api.GetDeployments)
		deployments.GET("/:deployment", perm.MustRead(), api.GetDeployment)
	} // end of deployments endpoints
}
