// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/repo"
)

// swagger:operation GET /api/v1/repos/{org}/{repo} repos GetRepo
//
// Get a repository
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

// GetRepo represents the API handler to get a repository.
func GetRepo(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)

	l.Debug("reading repo")

	c.JSON(http.StatusOK, r)
}
