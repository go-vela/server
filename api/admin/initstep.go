// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/admin/initstep admin AdminUpdateInitStep
//
// Update an initstep in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing initstep to update
//   required: true
//   schema:
//     "$ref": "#/definitions/InitStep"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the initstep in the database
//     schema:
//       "$ref": "#/definitions/InitStep"
//   '404':
//     description: Unable to update the initstep in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the initstep in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateInitStep represents the API handler to
// update any initstep stored in the database.
func UpdateInitStep(c *gin.Context) {
	logrus.Info("Admin: updating initstep in database")

	// capture body from API request
	input := new(library.InitStep)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for initstep %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the initstep
	err = database.FromContext(c).UpdateInitStep(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update initstep %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
