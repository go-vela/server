// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
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

// swagger:operation PUT /api/v1/admin/service admin AdminUpdateService
//
// Update a hook in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing service to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Service"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the service in the database
//     type: json
//     schema:
//       "$ref": "#/definitions/Service"
//   '404':
//     description: Unable to update the service in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the service in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateService represents the API handler to
// update any service stored in the database.
func UpdateService(c *gin.Context) {
	logrus.Info("Admin: updating service in database")

	// capture middleware values
	ctx := c.Request.Context()

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

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
