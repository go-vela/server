// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/dashboard"
	dMiddleware "github.com/go-vela/server/router/middleware/dashboard"
)

// DashboardHandlers is a function that extends the provided base router group
// with the API handlers for dashboard functionality.
//
// POST   /api/v1/dashboards
// GET    /api/v1/dashboards/:id
// PUT    /api/v1/dashboards/:id
// DELETE /api/v1/dashboards/:id .
func DashboardHandlers(base *gin.RouterGroup) {
	// Dashboard endpoints
	dashboards := base.Group("/dashboards")
	{
		dashboards.POST("", dashboard.CreateDashboard)

		d := dashboards.Group("/:dashboard", dMiddleware.Establish())
		{
			d.GET("", dashboard.GetDashboard)
			d.PUT("", dashboard.UpdateDashboard)
			d.DELETE("", dashboard.DeleteDashboard)
		}
	} // end of dashboard endpoints
}
