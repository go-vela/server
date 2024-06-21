// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/hook"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/hooks/{org}/{repo}/{hook} webhook DeleteHook
//
// Delete a hook
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
//     description: Successfully deleted the webhook
//     schema:
//       type: string
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

// DeleteHook represents the API handler to remove a webhook.
func DeleteHook(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	h := hook.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), h.GetNumber())

	l.Debugf("deleting hook %s", entry)

	// send API call to remove the webhook
	err := database.FromContext(c).DeleteHook(ctx, h)
	if err != nil {
		retErr := fmt.Errorf("unable to delete hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("hook %s deleted", entry))
}
