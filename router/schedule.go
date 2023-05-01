// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/schedules"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/repo"
	smiddleware "github.com/go-vela/server/router/middleware/schedule"
)

// ScheduleHandler is a function that extends the provided base router group
// with the API handlers for schedule functionality.
//
// POST   /api/v1/schedules/:org/:repo
// GET    /api/v1/schedules/:org/:repo
// GET    /api/v1/schedules/:org/:repo/:schedule
// PUT    /api/v1/schedules/:org/:repo/:schedule
// DELETE /api/v1/schedules/:org/:repo/:schedule .
func ScheduleHandler(base *gin.RouterGroup) {
	// Schedules endpoints
	_schedules := base.Group("/schedules/:org/:repo", org.Establish(), repo.Establish(), smiddleware.Establish())
	{
		_schedules.POST("", perm.MustAdmin(), middleware.Payload(), schedules.CreateSchedule)
		_schedules.GET("", perm.MustRead(), schedules.ListSchedules)
		_schedules.GET("/:schedule", perm.MustRead(), schedules.GetSchedule)
		_schedules.PUT("/:schedule", perm.MustAdmin(), middleware.Payload(), schedules.UpdateSchedule)
		_schedules.DELETE("/:schedule", perm.MustAdmin(), schedules.DeleteSchedule)
	} // end of schedules endpoints
}
