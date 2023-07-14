// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with service
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/step"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/perm"
	smiddleware "github.com/go-vela/server/router/middleware/step"
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
// DELETE /api/v1/repos/:org/:repo/builds/:build/steps/:step/logs .
func StepHandlers(base *gin.RouterGroup) {
	// Steps endpoints
	steps := base.Group("/steps")
	{
		steps.POST("", perm.MustPlatformAdmin(), middleware.Payload(), step.CreateStep)
		steps.GET("", perm.MustRead(), step.ListSteps)

		// Step endpoints
		s := steps.Group("/:step", smiddleware.Establish())
		{
			s.GET("", perm.MustRead(), step.GetStep)
			s.PUT("", perm.MustBuildAccess(), middleware.Payload(), step.UpdateStep)
			s.DELETE("", perm.MustPlatformAdmin(), step.DeleteStep)

			// Log endpoints
			LogStepHandlers(s)
		} // end of step endpoints
	} // end of steps endpoints
}
