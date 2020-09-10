// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this orgsitory.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/org"
)

// orgHandlers is a function that extends the provided base router group
// with the API handlers for org functionality.
//
// POST   /api/v1/orgs
// GET    /api/v1/orgs
// GET    /api/v1/orgs/:org/:org
// PUT    /api/v1/orgs/:org/:org
// DELETE /api/v1/orgs/:org/:org
// PATCH  /api/v1/orgs/:org/:org/repair
// PATCH  /api/v1/orgs/:org/:org/chown
// POST   /api/v1/orgs/:org/:org/builds
// GET    /api/v1/orgs/:org/:org/builds
// POST   /api/v1/orgs/:org/:org/builds/:build
// GET    /api/v1/orgs/:org/:org/builds/:build
// PUT    /api/v1/orgs/:org/:org/builds/:build
// DELETE /api/v1/orgs/:org/:org/builds/:build
// GET    /api/v1/orgs/:org/:org/builds/:build/logs
// POST   /api/v1/orgs/:org/:org/builds/:build/services
// GET    /api/v1/orgs/:org/:org/builds/:build/services
// GET    /api/v1/orgs/:org/:org/builds/:build/services/:service
// PUT    /api/v1/orgs/:org/:org/builds/:build/services/:service
// DELETE /api/v1/orgs/:org/:org/builds/:build/services/:service
// POST   /api/v1/orgs/:org/:org/builds/:build/services/:service/logs
// GET    /api/v1/orgs/:org/:org/builds/:build/services/:service/logs
// PUT    /api/v1/orgs/:org/:org/builds/:build/services/:service/logs
// DELETE /api/v1/orgs/:org/:org/builds/:build/services/:service/logs
// POST   /api/v1/orgs/:org/:org/builds/:build/steps
// GET    /api/v1/orgs/:org/:org/builds/:build/steps
// GET    /api/v1/orgs/:org/:org/builds/:build/steps/:step
// PUT    /api/v1/orgs/:org/:org/builds/:build/steps/:step
// DELETE /api/v1/orgs/:org/:org/builds/:build/steps/:step
// POST   /api/v1/orgs/:org/:org/builds/:build/steps/:step/logs
// GET    /api/v1/orgs/:org/:org/builds/:build/steps/:step/logs
// PUT    /api/v1/orgs/:org/:org/builds/:build/steps/:step/logs
// DELETE /api/v1/orgs/:org/:org/builds/:build/steps/:step/logs
func OrgHandlers(base *gin.RouterGroup) {
	// org endpoint.
	orgs := base.Group("/org")
	{
		org := orgs.Group("/:org", org.Establish())
		{
			org.GET("", api.GetBuildOrgs)
		} // end of org endpoints
	} // end of orgs endpoints
}
