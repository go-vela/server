// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/user/token users CreateToken
//
// Create a token for the current authenticated user
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully created a token for the current user
//     schema:
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to create a token for the current user
//     schema:
//       "$ref": "#/definitions/Error"

// CreateToken represents the API handler to create
// a user token.
func CreateToken(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("composing token for user %s", u.GetName())

	tm := c.MustGet("token-manager").(*token.Manager)

	// compose JWT token for user
	rt, at, err := tm.Compose(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	u.SetRefreshToken(rt)

	// send API call to update the user
	_, err = database.FromContext(c).UpdateUser(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	l.Info("user updated - token created")

	c.JSON(http.StatusOK, types.Token{Token: &at})
}
