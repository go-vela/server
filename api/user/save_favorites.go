// SPDX-License-Identifier: Apache-2.0

package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/user/favorites users SaveUserFavorites
//
// Save the current authenticated user's favorites
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '204':
//     description: Successfully saved the current user's favorites
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// SaveUserFavorites represents the API handler to save the
// currently authenticated user's favorites.
func SaveUserFavorites(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	favorites := new([]*types.Favorite)

	err := c.Bind(favorites)
	if err != nil {
		retErr := err

		util.HandleError(ctx, http.StatusBadRequest, retErr)

		return
	}

	err = database.FromContext(c).UpdateFavorites(ctx, u, *favorites)
	if err != nil {
		retErr := err

		util.HandleError(ctx, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, favorites)
}
