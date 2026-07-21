// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/favorite"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
)

// BuildHandlers is a function that extends the provided base router group
// with the API handlers for build functionality.
//
// POST   /api/v1/repos/:org/:repo/builds
// GET    /api/v1/repos/:org/:repo/builds
// POST   /api/v1/repos/:org/:repo/builds/:build
// GET    /api/v1/repos/:org/:repo/builds/:build .
func FavoriteHandlers(base *gin.RouterGroup) {
	// Favorites endpoints
	favorites := base.Group("/favorites")
	{
		favorites.POST("", favorite.CreateFavorite)
		favorites.GET("", favorite.ListFavorites)

		// Favorite endpoints
		f := favorites.Group("/:org/:repo", org.Establish(), repo.Establish())
		{
			f.DELETE("", favorite.DeleteFavorite)
			f.PUT("", favorite.UpdateFavorite)
		} // end of favorite endpoints
	} // end of favorites endpoints
}
