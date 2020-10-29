// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/repo"
)

// PipelineHandlers is a function that extends the provided base router group
// with the API handlers for pipeline functionality.
//
// GET  /api/v1/pipelines/:org/:repo
// GET  /api/v1/pipelines/:org/:repo/templates
// POST /api/v1/pipelines/:org/:repo/expand
// POST /api/v1/pipelines/:org/:repo/compile
// POST /api/v1/pipelines/:org/:repo/validate
func PipelineHandlers(base *gin.RouterGroup) {
	// Pipelines endpoints
	pipelines := base.Group("pipelines/:org/:repo", repo.Establish())
	{
		pipelines.GET("", api.GetPipeline)
		pipelines.GET("/templates", api.GetTemplates)
		pipelines.POST("/expand", api.ExpandPipeline)
		pipelines.POST("/validate", api.ValidatePipeline)
		pipelines.POST("/compile", api.CompilePipeline)
	} // end of hooks endpoints
}
