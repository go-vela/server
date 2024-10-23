// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code to step
package log

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services CreateServiceLog
//
// Create the logs for a service
//
// ---
// deprecated: true
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
//   description: Payload containing the log to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Log"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the service logs
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

// CreateServiceLog represents the API handler to create
// the logs for a service.
func CreateServiceLog(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	l.Debugf("creating logs for service %s", entry)

	// capture body from API request
	input := new(types.Log)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in log object
	input.SetServiceID(s.GetID())
	input.SetBuildID(b.GetID())
	input.SetRepoID(r.GetID())

	// send API call to create the logs
	err = database.FromContext(c).CreateLog(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"service":    s.GetName(),
		"service_id": s.GetID(),
	}).Info("logs created for service")

	c.JSON(http.StatusCreated, nil)
}
