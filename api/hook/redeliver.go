// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/hook"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/settings"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/hooks/{org}/{repo}/{hook}/redeliver webhook RedeliverHook
//
// Redeliver a hook
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
//     description: Successfully redelivered the webhook
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

// RedeliverHook represents the API handler to redeliver
// a webhook from the SCM.
func RedeliverHook(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	h := hook.Retrieve(c)
	ps := settings.FromContext(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), h.GetNumber())

	// check to see if queue has reached configured capacity to allow for a redelivery (hook must have a build)
	if h.GetBuild().GetStatus() == constants.StatusPending && ps.GetQueueRestartLimit() > 0 {
		// check length of specified route for the build
		queueLength, err := queue.FromContext(c).RouteLength(ctx, h.GetBuild().GetRoute())
		if err != nil {
			util.HandleError(c, http.StatusInternalServerError, fmt.Errorf("unable to get queue length for %s: %w", h.GetBuild().GetRoute(), err))

			return
		}

		if queueLength >= int64(ps.GetQueueRestartLimit()) {
			retErr := fmt.Errorf("unable to restart build %s: queue length %d exceeds configured limit %d, please wait for the queue to decrease in size before retrying", entry, queueLength, ps.GetQueueRestartLimit())

			util.HandleError(c, http.StatusTooManyRequests, retErr)

			return
		}
	}

	l.Debugf("redelivering hook %s", entry)

	err := scm.FromContext(c).RedeliverWebhook(c, u, h)
	if err != nil {
		retErr := fmt.Errorf("unable to redeliver hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("hook %s redelivered", entry))
}
