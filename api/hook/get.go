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

// swagger:operation GET /api/v1/hooks/{org}/{repo}/{hook} webhook GetHook
//
// Retrieve a webhook for the configured backend
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
//     description: Successfully retrieved the webhook
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: Unable to retrieve the webhook
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the webhook
//     schema:
//       "$ref": "#/definitions/Error"

// GetHook represents the API handler to capture a
// webhook from the configured backend.
func GetHook(c *gin.Context) {
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
	}).Infof("reading hook %s", entry)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHookForRepo(ctx, r, number)
	if err != nil {
		retErr := fmt.Errorf("unable to get hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, h)
}
