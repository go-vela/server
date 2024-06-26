// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/repos/{org}/{repo} repos DeleteRepo
//
// Delete a repository
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
//     description: Successfully deleted the repo
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

// DeleteRepo represents the API handler to remove a repository.
func DeleteRepo(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("deleting repo %s", r.GetFullName())

	// send API call to remove the webhook
	err := scm.FromContext(c).Disable(ctx, u, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to delete webhook for %s: %w", r.GetFullName(), err)

		if err.Error() == "Repo not found" {
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Mark the repo as inactive
	r.SetActive(false)

	_, err = database.FromContext(c).UpdateRepo(ctx, r)
	if err != nil {
		retErr := fmt.Errorf("unable to set repo %s to inactive: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Comment out actual delete until delete mechanism is fleshed out
	// err = database.FromContext(c).DeleteRepo(r.ID)
	// if err != nil {
	// 	retErr := fmt.Errorf("Error while deleting repo %s: %w", r.FullName, err)
	// 	util.HandleError(c, http.StatusInternalServerError, retErr)
	// 	return
	// }

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s set to inactive", r.GetFullName()))
}
