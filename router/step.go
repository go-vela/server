// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// nolint: dupl // ignore similar code with service
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/step"
)

// StepHandlers is a function that extends the provided base router group
// with the API handlers for step functionality.
//
// POST   /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step
// POST   /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// GET    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// PUT    /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs
// POST   /api/v1/repos/:org/:repo/builds/:build/steps/:step/stream .
func StepHandlers(base *gin.RouterGroup) {
	// Steps endpoints
	steps := base.Group("/steps")
	{
		steps.POST("", perm.MustPlatformAdmin(), middleware.Payload(), api.CreateStep)
		steps.GET("", perm.MustRead(), api.GetSteps)

		// Step endpoints
		step := steps.Group("/:step", step.Establish())
		{
			step.GET("", perm.MustRead(), api.GetStep)
			step.PUT("", perm.MustPlatformAdmin(), middleware.Payload(), api.UpdateStep)
			step.DELETE("", perm.MustPlatformAdmin(), api.DeleteStep)

			step.POST("/stream", api.PostStepStream)

			// Log endpoints
			LogStepHandlers(step)
		} // end of step endpoints
	} // end of steps endpoints
}
