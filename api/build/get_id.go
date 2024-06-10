// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/search/builds/{id} builds GetBuildByID
//
// Get a build by id
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: id
//   description: Build ID
//   required: true
//   type: number
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved build
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildByID represents the API handler to get a
// build by its id.
func GetBuildByID(c *gin.Context) {
	// Capture user from middleware
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// Parse build ID from path
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to parse build id: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": id,
		"user":  u.GetName(),
	}).Debugf("reading build %d", id)

	// Get build from database
	b, err := database.FromContext(c).GetBuild(ctx, id)
	if err != nil {
		retErr := fmt.Errorf("unable to get build: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Capture user access from SCM. We do this in order to ensure user has access and is not
	// just retrieving any build using a random id number.
	perm, err := scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), b.GetRepo().GetOrg(), b.GetRepo().GetName())
	if err != nil {
		logrus.Errorf("unable to get user %s access level for repo %s", u.GetName(), b.GetRepo().GetFullName())
	}

	// Ensure that user has at least read access to repo to return the build
	if perm == "none" && !u.GetAdmin() {
		retErr := fmt.Errorf("unable to retrieve build %d: user does not have read access to repo", id)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}
