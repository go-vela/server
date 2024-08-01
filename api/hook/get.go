// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/hook"
	"github.com/go-vela/server/router/middleware/repo"
)

// swagger:operation GET /api/v1/hooks/{org}/{repo}/{hook} webhook GetHook
//
// Get a hook
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
// - in: path
//   name: hook
//   description: Number of the hook
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the webhook
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetHook represents the API handler to get a hook.
func GetHook(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	h := hook.Retrieve(c)

	l.Debugf("reading hook %s/%d", r.GetFullName(), h.GetNumber())

	c.JSON(http.StatusOK, h)
}
