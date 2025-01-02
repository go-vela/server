// SPDX-License-Identifier: Apache-2.0

package build

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/build"
)

// swagger:operation GET /status/{org}/{repo}/{build} builds GetStatus
//
// Get a build status
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the build
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Build"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Build"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Build"

// GetStatus represents the API handler to return "status", a lite representation of the resource with limited fields for unauthenticated access.
func GetStatus(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)

	l.Debug("reading status for build")

	// sanitize fields for the unauthenticated response
	b.StatusSanitize()

	c.JSON(http.StatusOK, b)
}
