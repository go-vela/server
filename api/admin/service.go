// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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

// swagger:operation GET /api/v1/admin/services admin AdminAllServices
//
// Get all of the services in the database
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved all services from the database
//     type: json
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Service"
//   '500':
//     description: Unable to retrieve all services from the database
//     schema:
//       type: string

// AllServices represents the API handler to
// captures all services stored in the database.
func AllServices(c *gin.Context) {
	logrus.Info("Admin: reading all services")

	// send API call to capture all services
	s, err := database.FromContext(c).GetServiceList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all services: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}

// swagger:operation PUT /api/v1/admin/service admin AdminUpdateService
//
// Update a hook in the database
//
// ---
// x-success_http_code: '200'
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
//       type: string
//   '501':
//     description: Unable to update the service in the database
//     schema:
//       type: string

// UpdateService represents the API handler to
// update any service stored in the database.
func UpdateService(c *gin.Context) {
	logrus.Info("Admin: updating service in database")

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the service
	err = database.FromContext(c).UpdateService(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
