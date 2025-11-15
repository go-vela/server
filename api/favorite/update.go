// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/user/favorites/{org}/{repo} favorites UpdateFavorite
//
// Update the current authenticated user's favorite
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '204':
//     description: Successfully updated favorite position
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateFavorite represents the API handler to update the
// currently authenticated user's favorite position.
func UpdateFavorite(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	favorite := new(types.Favorite)

	err := c.Bind(favorite)
	if err != nil {
		retErr := err

		util.HandleError(ctx, http.StatusBadRequest, retErr)

		return
	}

	err = database.FromContext(c).UpdateFavoritePosition(ctx, u, r, favorite)
	if err != nil {
		retErr := fmt.Errorf("unable to update favorite position for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.Status(http.StatusNoContent)
}
