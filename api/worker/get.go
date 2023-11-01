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

// swagger:operation GET /api/v1/workers/{worker} workers GetWorker
//
// Retrieve a worker for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: Hostname of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the worker
//     schema:
//       "$ref": "#/definitions/Worker"
//   '404':
//     description: Unable to retrieve the worker
//     schema:
//       "$ref": "#/definitions/Error"

// GetWorker represents the API handler to capture a
// worker from the configured backend.
func GetWorker(c *gin.Context) {
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
	}).Infof("reading worker %s", w.GetHostname())

	var rBs []*library.Build
	for _, b := range w.GetRunningBuilds() {
		build, err := database.FromContext(c).GetBuild(ctx, b.GetID())
		if err != nil {
			retErr := fmt.Errorf("unable to read build %d: %w", b.GetID(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		rBs = append(rBs, build)
	}

	w.SetRunningBuilds(rBs)

	c.JSON(http.StatusOK, w)
}
