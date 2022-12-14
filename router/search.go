// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
)

// SearchHandlers is a function that extends the provided base router group
// with the API handlers for resource search functionality.
//
// GET    /api/v1/search/builds/:id .
func SearchHandlers(base *gin.RouterGroup) {
	// Search endpoints
	search := base.Group("/search")
	{
		// Build endpoint
		build := search.Group("/builds")
		{
			build.GET("/:id", api.GetBuildByID)
		}
	} // end of search endpoints
}
