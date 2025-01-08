// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/repo"
)

// swagger:operation GET /status/{org}/{repo} repos GetRepoStatus
//
// Get a repository status
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     schema:
//       "$ref": "#/definitions/Repo"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Repo"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Repo"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Repo"

// GetRepoStatus represents the API handler to return "status", a lite representation of the resource with limited fields for unauthenticated access.
func GetRepoStatus(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)

	l.Debug("reading status for repo")

	// sanitize fields for the unauthenticated response
	r.StatusSanitize()

	c.JSON(http.StatusOK, r)
}
