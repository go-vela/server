// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/user/favorites favorites CreateFavorite
//
// Save the current authenticated user's favorites
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully added user favorite
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// CreateFavorite represents the API handler to add a
// favorite for the currently authenticated user.
func CreateFavorite(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	favorite := new(types.Favorite)

	err := c.Bind(favorite)
	if err != nil {
		retErr := err

		util.HandleError(ctx, http.StatusBadRequest, retErr)

		return
	}

	err = database.FromContext(c).CreateFavorite(ctx, u, favorite)
	if err != nil {
		retErr := fmt.Errorf("unable to create favorite for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.Status(http.StatusCreated)
}
