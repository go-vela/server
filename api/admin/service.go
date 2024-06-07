// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation PUT /api/v1/admin/service admin AdminUpdateService
//
// Update a service
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The service object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Service"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the service
//     type: json
//     schema:
//       "$ref": "#/definitions/Service"
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

// UpdateService represents the API handler to update a service.
func UpdateService(c *gin.Context) {
	logrus.Info("Admin: updating service in database")

	// capture middleware values
	ctx := c.Request.Context()

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to update the service
	s, err := database.FromContext(c).UpdateService(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
