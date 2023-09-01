// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"fmt"
	"net/http"

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

	w, err := database.FromContext(c).ListWorkers(ctx)
	if err != nil {
		retErr := fmt.Errorf("unable to get workers: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, w)
}
