// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/perm"
)

// DeploymentHandlers is a function that extends the provided base router group
// with the API handlers for deployment functionality.
//
// POST   /api/v1/deployments/:org/:repo
// GET    /api/v1/deployments/:org/:repo
// GET    /api/v1/deployments/:org/:repo/:deployment
// PUT    /api/v1/deployments/:org/:repo/:deployment
// DELETE /api/v1/deployments/:org/:repo/:deployment
func DeploymentHandlers(base *gin.RouterGroup) {
	// Deployments endpoints
	deployments := base.Group("/deployments/:org/:repo")
	{
		deployments.POST("", perm.MustAdmin(), api.CreateDeployment)
		deployments.GET("", perm.MustAdmin(), api.GetDeployments)
		deployments.GET("/:deployment", perm.MustAdmin(), api.GetDeployment)
		deployments.PUT("/:deployment", perm.MustAdmin(), api.UpdateDeployment)
		deployments.DELETE("/:deployment", perm.MustAdmin(), api.DeleteDeployment)
	} // end of deployments endpoints
}
