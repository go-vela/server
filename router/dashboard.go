// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/dashboard"
	dMiddleware "github.com/go-vela/server/router/middleware/dashboard"
)

// DashboardHandlers is a function that extends the provided base router group
// with the API handlers for resource search functionality.
//
// GET    /api/v1/search/builds/:id .
func DashboardHandlers(base *gin.RouterGroup) {
	// Search endpoints
	dashboards := base.Group("/dashboards")
	{
		dashboards.POST("", dashboard.CreateDashboard)

		d := dashboards.Group("/:dashboard", dMiddleware.Establish())
		{
			d.GET("", dashboard.GetDashboard)
			d.PUT("", dashboard.UpdateDashboard)
		}
	} // end of search endpoints
}
