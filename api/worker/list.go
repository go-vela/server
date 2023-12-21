// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/workers workers ListWorkers
//
// Retrieve a list of workers for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: active
//   description: Filter workers based on active status
//   type: string
//   enum:
//   - true
//   - false
//   - all
// - in: query
//   name: checked_in_before
//   description: filter workers that have checked in before a certain time
//   type: integer
// - in: query
//   name: checked_in_after
//   description: filter workers that have checked in after a certain time
//   type: integer
//   default: 0
// security:
//   - ApiKeyAuth: []
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

	active := c.DefaultQuery("active", "true")

	// capture before query parameter if present, default to now
	before, err := strconv.ParseInt(c.DefaultQuery("checked_in_before", strconv.FormatInt(time.Now().UTC().Unix(), 10)), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert `checked_in_before` query parameter: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture after query parameter if present, default to 0
	after, err := strconv.ParseInt(c.DefaultQuery("checked_in_after", "0"), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert `checked_in_after` query parameter: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	w, err := database.FromContext(c).ListWorkers(ctx, active, before, after)
	if err != nil {
		retErr := fmt.Errorf("unable to get workers: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, w)
}
