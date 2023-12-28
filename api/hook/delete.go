// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/hooks/{org}/{repo}/{hook} webhook DeleteHook
//
// Delete a webhook for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
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
//     description: The webhook was unable to be deleted
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: The webhook was unable to be deleted
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: The webhook was unable to be deleted
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteHook represents the API handler to remove
// a webhook from the configured backend.
func DeleteHook(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	hook := util.PathParameter(c, "hook")
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), hook)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"hook": hook,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("deleting hook %s", entry)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHookForRepo(ctx, r, number)
	if err != nil {
		retErr := fmt.Errorf("unable to get hook %s: %w", hook, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to remove the webhook
	err = database.FromContext(c).DeleteHook(ctx, h)
	if err != nil {
		retErr := fmt.Errorf("unable to delete hook %s: %w", hook, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("hook %s deleted", entry))
}
