// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/workers workers ListWorkers
//
// Retrieve a list of workers for the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
// - in: query
//   name: links
//   description: Include links to builds currently running
//   type: boolean
//   default: false
// responses:
//   '200':
//     description: Successfully retrieved the list of workers
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Worker"
//   '500':
//     description: Unable to retrieve the list of workers
//     schema:
//       "$ref": "#/definitions/Error"

// ListWorkers represents the API handler to capture a
// list of workers from the configured backend.
func ListWorkers(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Info("reading workers")

	var filters = map[string]interface{}{}

	active := c.Query("active")

	if len(active) > 0 {
		filters["active"] = active
	}

	workers, err := database.FromContext(c).ListWorkers(ctx)
	if err != nil {
		retErr := fmt.Errorf("unable to get workers: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	for _, w := range workers {
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
	}

	c.JSON(http.StatusOK, workers)
}
