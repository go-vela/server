// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with user.go
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation PUT /api/v1/admin/step admin AdminUpdateStep
//
// Update a step
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The step object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Step"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the step
//     schema:
//       "$ref": "#/definitions/Step"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateStep represents the API handler to update a step.
func UpdateStep(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()
	u := user.Retrieve(c)

	logger := logrus.WithFields(logrus.Fields{
		"ip":      util.EscapeValue(c.ClientIP()),
		"path":    util.EscapeValue(c.Request.URL.Path),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	})

	logrus.Debug("platform admin: updating step")

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"step_id": input.GetID(),
		"step":    util.EscapeValue(input.GetName()),
	}).Debug("attempting to update step")

	// send API call to update the step
	s, err := database.FromContext(c).UpdateStep(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update step %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"step_id": s.GetID(),
		"step":    s.GetName(),
	}).Info("updated step")

	c.JSON(http.StatusOK, s)
}
