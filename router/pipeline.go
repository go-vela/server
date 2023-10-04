// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/pipeline"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/perm"
	pmiddleware "github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
)

// PipelineHandlers is a function that extends the provided base router group
// with the API handlers for pipeline functionality.
//
// POST   /api/v1/pipelines/:org/:repo
// GET    /api/v1/pipelines/:org/:repo
// GET    /api/v1/pipelines/:org/:repo/:pipeline
// PUT    /api/v1/pipelines/:org/:repo/:pipeline
// DELETE /api/v1/pipelines/:org/:repo/:pipeline
// GET    /api/v1/pipelines/:org/:repo/:pipeline/templates
// POST   /api/v1/pipelines/:org/:repo/:pipeline/expand
// POST   /api/v1/pipelines/:org/:repo/:pipeline/compile
// POST   /api/v1/pipelines/:org/:repo/:pipeline/validate .
func PipelineHandlers(base *gin.RouterGroup) {
	// Pipelines endpoints
	_pipelines := base.Group("pipelines/:org/:repo", org.Establish(), repo.Establish())
	{
		_pipelines.POST("", perm.MustAdmin(), pipeline.CreatePipeline)
		_pipelines.GET("", perm.MustRead(), pipeline.ListPipelines)

		_pipeline := _pipelines.Group("/:pipeline", pmiddleware.Establish())
		{
			_pipeline.GET("", perm.MustRead(), pipeline.GetPipeline)
			_pipeline.PUT("", perm.MustWrite(), pipeline.UpdatePipeline)
			_pipeline.DELETE("", perm.MustPlatformAdmin(), pipeline.DeletePipeline)
			_pipeline.GET("/templates", perm.MustRead(), pipeline.GetTemplates)
			_pipeline.POST("/compile", perm.MustRead(), pipeline.CompilePipeline)
			_pipeline.POST("/expand", perm.MustRead(), pipeline.ExpandPipeline)
			_pipeline.POST("/validate", perm.MustRead(), pipeline.ValidatePipeline)
		} // end of pipeline endpoints
	} // end of pipelines endpoints
}
