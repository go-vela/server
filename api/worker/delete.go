// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/workers/{worker} workers DeleteWorker
//
// Delete a worker
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: Name of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted worker
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

// DeleteWorker represents the API handler to remove a worker.
func DeleteWorker(c *gin.Context) {
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
	}).Infof("deleting worker %s", w.GetHostname())

	// send API call to remove the step
	err := database.FromContext(c).DeleteWorker(ctx, w)
	if err != nil {
		retErr := fmt.Errorf("unable to delete worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("worker %s deleted", w.GetHostname()))
}
