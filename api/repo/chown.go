// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/chown repos ChownRepo
//
// Change the owner of a repository
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
//     description: Successfully changed the owner for the repository
//     schema:
//       type: string
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

// ChownRepo represents the API handler to change
// the owner of a repo.
func ChownRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"ip":      util.EscapeValue(c.ClientIP()),
		"path":    util.EscapeValue(c.Request.URL.Path),
		"org":     o,
		"repo":    r.GetName(),
		"repo_id": r.GetID(),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	})

	logger.Debugf("changing owner of repo %s to %s", r.GetFullName(), u.GetName())

	// update repo owner
	r.SetOwner(u)

	// send API call to update the repo
	_, err := database.FromContext(c).UpdateRepo(ctx, r)
	if err != nil {
		retErr := fmt.Errorf("unable to change owner of repo %s to %s: %w", r.GetFullName(), u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logger.Infof("updated repo - changed owner to %s", u.GetName())

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s changed owner to %s", r.GetFullName(), u.GetName()))
}
