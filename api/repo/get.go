// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
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
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("reading repo %s", r.GetFullName())

	c.JSON(http.StatusOK, r)
}
