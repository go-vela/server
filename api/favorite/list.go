// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/user/favorites favorites ListFavorites
//
// Get the current authenticated user's favorites
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the current user's favorites
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Favorite"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// ListFavorites represents the API handler to capture the
// currently authenticated user's favorites.
func ListFavorites(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	favorites, err := database.FromContext(c).ListUserFavorites(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to get favorites for user %s: %w", u.GetName(), err)

		util.HandleError(ctx, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, favorites)
}
