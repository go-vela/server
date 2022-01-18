// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// nolint: dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/admin/steps admin AdminAllSteps
//
// Get all of the steps in the database
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved all steps from the database
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Step"
//   '500':
//     description: Unable to retrieve all steps from the database
//     schema:
//       "$ref": "#/definitions/Error"

// AllSteps represents the API handler to
// captures all steps stored in the database.
func AllSteps(c *gin.Context) {
	logrus.Info("Admin: reading all steps")

	// send API call to capture all steps
	s, err := database.FromContext(c).GetStepList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all steps: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}

// swagger:operation PUT /api/v1/admin/step admin AdminUpdateStep
//
// Update a step in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing step to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Step"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the step in the database
//     schema:
//       "$ref": "#/definitions/Step"
//   '404':
//     description: Unable to update the step in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the step in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateStep represents the API handler to
// update any step stored in the database.
func UpdateStep(c *gin.Context) {
	logrus.Info("Admin: updating step in database")

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the step
	err = database.FromContext(c).UpdateStep(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update step %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
