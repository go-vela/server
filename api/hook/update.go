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
	"github.com/go-vela/types/library"
)

// swagger:operation PUT /api/v1/hooks/{org}/{repo}/{hook} webhook UpdateHook
//
// Update a hook
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
// - in: body
//   name: body
//   description: Webhook payload that we expect from the user or VCS
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the webhook
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

// UpdateHook represents the API handler to update a hook.
func UpdateHook(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	h := hook.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), h.GetNumber())

	l.Debugf("updating hook %s", entry)

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for hook %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update webhook fields if provided
	if input.GetCreated() > 0 {
		// update created if set
		h.SetCreated(input.GetCreated())
	}

	if len(input.GetHost()) > 0 {
		// update host if set
		h.SetHost(input.GetHost())
	}

	if len(input.GetEvent()) > 0 {
		// update event if set
		h.SetEvent(input.GetEvent())
	}

	if len(input.GetBranch()) > 0 {
		// update branch if set
		h.SetBranch(input.GetBranch())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		h.SetError(input.GetError())
	}

	if len(input.GetStatus()) > 0 {
		// update status if set
		h.SetStatus(input.GetStatus())
	}

	if len(input.GetLink()) > 0 {
		// update link if set
		h.SetLink(input.GetLink())
	}

	// send API call to update the webhook
	h, err = database.FromContext(c).UpdateHook(ctx, h)
	if err != nil {
		retErr := fmt.Errorf("unable to update hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, h)
}
