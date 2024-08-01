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

// swagger:operation GET /api/v1/workers/{worker} workers GetWorker
//
// Get a worker
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

// GetWorker represents the API handler to get a worker.
func GetWorker(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	w := worker.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("reading worker %s", w.GetHostname())

	rBs := []*types.Build{}

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
