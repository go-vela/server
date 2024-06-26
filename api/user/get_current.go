// SPDX-License-Identifier: Apache-2.0

package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/user"
)

// swagger:operation GET /api/v1/user users GetCurrentUser
//
// Get the current authenticated user
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the current user
//     schema:
//       "$ref": "#/definitions/User"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// GetCurrentUser represents the API handler to capture the
// currently authenticated user.
func GetCurrentUser(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)

	l.Debugf("reading current user %s", u.GetName())

	c.JSON(http.StatusOK, u)
}
