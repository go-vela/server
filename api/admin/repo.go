// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with service.go
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/admin/repo admin AdminUpdateRepo
//
// Update a repository
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The repository object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the repo
//     schema:
//       "$ref": "#/definitions/Repo"
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

// UpdateRepo represents the API handler to update a repo.
func UpdateRepo(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()
	u := user.Retrieve(c)

	logger := logrus.WithFields(logrus.Fields{
		"ip":      util.EscapeValue(c.ClientIP()),
		"path":    util.EscapeValue(c.Request.URL.Path),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	})

	logrus.Debug("platform admin: updating repo")

	// capture body from API request
	input := new(types.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for repo %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"repo_id": input.GetID(),
		"repo":    util.EscapeValue(input.GetFullName()),
	}).Debug("attempting to update repo")

	// send API call to update the repo
	r, err := database.FromContext(c).UpdateRepo(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update repo %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"repo_id": r.GetID(),
		"repo":    r.GetFullName(),
	}).Info("repo updated")

	c.JSON(http.StatusOK, r)
}
