// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/artifact"
	tamiddleware "github.com/go-vela/server/router/middleware/artifact"
)

// ArtifactHandlers is a function that extends the provided base router group
// with the API handlers for artifact functionality.
//
// GET    /api/v1/repos/:org/:repo/builds/:build/artifact
// GET    /api/v1/repos/:org/:repo/builds/:build/artifact/:artifact
// PUT    /api/v1/repos/:org/:repo/builds/:build/artifact .
func ArtifactHandlers(base *gin.RouterGroup) {
	// artifact endpoints
	_artifact := base.Group("/artifact")
	{
		_artifact.GET("", artifact.ListArtifactsForBuild)
		_artifact.PUT("", artifact.CreateArtifact)

		// Individual artifact endpoints
		a := _artifact.Group("/:artifact", tamiddleware.Establish())
		{
			a.GET("", artifact.GetArtifact)
		} // end of individual artifact endpoints
	} // end of artifact endpoints
}
