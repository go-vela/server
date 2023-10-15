// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/workers/{worker} workers UpdateWorker
//
// Update a worker for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the worker to update
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
//     description: Unable to update the worker
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to update the worker
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the worker
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateWorker represents the API handler to
// update a worker in the configured backend.
func UpdateWorker(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	w := worker.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user":   u.GetName(),
		"worker": w.GetHostname(),
	}).Infof("updating worker %s", w.GetHostname())

	// capture body from API request
	input := new(library.Worker)

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

	if input.RunningBuildIDs != nil {
		// update runningBuildIDs if set
		w.SetRunningBuildIDs(input.GetRunningBuildIDs())
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
