// SPDX-License-Identifier: Apache-2.0

package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

//
// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build}/services/{service} services UpdateService
//
// Update a service for a build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: service
//   description: Service number
//   required: true
//   type: integer
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
//     schema:
//       "$ref": "#/definitions/Service"
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

// UpdateService represents the API handler to update
// a service for a build.
func UpdateService(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	l.Debugf("updating service %s", entry)

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update service fields if provided
	if len(input.GetStatus()) > 0 {
		// update status if set
		s.SetStatus(input.GetStatus())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		s.SetError(input.GetError())
	}

	if input.GetExitCode() > 0 {
		// update exit_code if set
		s.SetExitCode(input.GetExitCode())
	}

	if input.GetStarted() > 0 {
		// update started if set
		s.SetStarted(input.GetStarted())
	}

	if input.GetFinished() > 0 {
		// update finished if set
		s.SetFinished(input.GetFinished())
	}

	// send API call to update the service
	s, err = database.FromContext(c).UpdateService(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to update service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
