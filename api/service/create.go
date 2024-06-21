// SPDX-License-Identifier: Apache-2.0

package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/services services CreateService
//
// Create a service for a build
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
// - in: body
//   name: body
//   description: Service object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Service"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the service
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

// CreateService represents the API handler to create
// a service for a build.
func CreateService(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	l.Debugf("creating new service for build %s", entry)

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new service for build %s: %w", entry, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in service object
	input.SetRepoID(r.GetID())
	input.SetBuildID(b.GetID())

	if len(input.GetStatus()) == 0 {
		input.SetStatus(constants.StatusPending)
	}

	if input.GetCreated() == 0 {
		input.SetCreated(time.Now().UTC().Unix())
	}

	// send API call to create the service
	s, err := database.FromContext(c).CreateService(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create service for build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"service":    s.GetName(),
		"service_id": s.GetID(),
	}).Infof("service %s created for build %s", s.GetName(), entry)

	c.JSON(http.StatusCreated, s)
}
