// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/schedule"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/perm"
	"github.com/go-vela/server/router/middleware/repo"
	sMiddleware "github.com/go-vela/server/router/middleware/schedule"
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
	_schedules := base.Group("/schedules/:org/:repo", org.Establish(), repo.Establish())
	{
		_schedules.POST("", perm.MustAdmin(), middleware.Payload(), schedule.CreateSchedule)
		_schedules.GET("", perm.MustRead(), schedule.ListSchedules)

		s := _schedules.Group("/:schedule", sMiddleware.Establish())
		{
			s.GET("", perm.MustRead(), schedule.GetSchedule)
			s.PUT("", perm.MustAdmin(), middleware.Payload(), schedule.UpdateSchedule)
			s.DELETE("", perm.MustAdmin(), schedule.DeleteSchedule)
		}
	} // end of schedules endpoints
}
