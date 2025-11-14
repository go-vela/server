// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/user/favorites/{org}/{repo} favorites DeleteFavorite
//
// Remove the current authenticated user's favorite
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '204':
//     description: Successfully removed user favorite
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteFavorite represents the API handler to delete a
// favorite for the currently authenticated user.
func DeleteFavorite(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	err := database.FromContext(c).DeleteFavorite(ctx, u, r)
	if err != nil {
		retErr := err

		util.HandleError(ctx, http.StatusInternalServerError, retErr)

		return
	}

	c.Status(http.StatusNoContent)
}
