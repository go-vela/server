// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/workers/{worker} workers UpdateWorker
//
// Update a worker
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The worker object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Worker"
// - in: path
//   name: worker
//   description: Name of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the worker
//     schema:
//       "$ref": "#/definitions/Worker"
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

// UpdateWorker represents the API handler to
// update a worker.
func UpdateWorker(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	w := worker.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("updating worker %s", w.GetHostname())

	// capture body from API request
	input := new(types.Worker)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if len(input.GetAddress()) > 0 {
		// update address if set
		w.SetAddress(input.GetAddress())
	}

	if len(input.GetRoutes()) > 0 {
		// update routes if set
		w.SetRoutes(input.GetRoutes())
	}

	if input.Active != nil {
		// update active if set
		w.SetActive(input.GetActive())
	}

	if input.RunningBuilds != nil {
		// update runningBuildIDs if set
		w.SetRunningBuilds(input.GetRunningBuilds())
	}

	if len(input.GetStatus()) > 0 {
		// update status if set
		w.SetStatus(input.GetStatus())
	}

	if input.GetLastStatusUpdateAt() > 0 {
		// update lastStatusUpdateAt if set
		w.SetLastStatusUpdateAt(input.GetLastStatusUpdateAt())
	}

	if input.GetLastBuildStartedAt() > 0 {
		// update lastBuildStartedAt if set
		w.SetLastBuildStartedAt(input.GetLastBuildStartedAt())
	}

	if input.GetLastBuildFinishedAt() > 0 {
		// update lastBuildFinishedAt if set
		w.SetLastBuildFinishedAt(input.GetLastBuildFinishedAt())
	}

	// send API call to update the worker
	w, err = database.FromContext(c).UpdateWorker(ctx, w)
	if err != nil {
		retErr := fmt.Errorf("unable to update worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, w)
}
