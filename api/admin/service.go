// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
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
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: updating service")

	// capture body from API request
	input := new(types.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"service_id": input.GetID(),
		"service":    util.EscapeValue(input.GetName()),
	}).Debug("platform admin: attempting to update service")

	// send API call to update the service
	s, err := database.FromContext(c).UpdateService(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"service_id": s.GetID(),
		"service":    s.GetName(),
	}).Info("platform admin: updated service")

	c.JSON(http.StatusOK, s)
}
