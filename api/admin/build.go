// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/admin/builds/queue admin AllBuildsQueue
//
// Get all of the running and pending builds in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: after
//   description: Unix timestamp to limit builds returned
//   required: false
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved all running and pending builds from the database
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/BuildQueue"
//   '401':
//     description: Unauthorized to retrieve running and pending builds
//     schema:
//       "$ref": "#/definitions/Error
//   '500':
//     description: Unable to retrieve all running and pending builds from the database
//     schema:
//       "$ref": "#/definitions/Error"

// AllBuildsQueue represents the API handler to
// capture all running and pending builds stored in the database.
func AllBuildsQueue(c *gin.Context) {
	logrus.Debug("platform admin: reading running and pending builds")

	// capture middleware values
	ctx := c.Request.Context()

	// default timestamp to 24 hours ago if user did not provide it as query parameter
	after := c.DefaultQuery("after", strconv.FormatInt(time.Now().UTC().Add(-24*time.Hour).Unix(), 10))

	// send API call to capture pending and running builds
	b, err := database.FromContext(c).ListPendingAndRunningBuilds(ctx, after)
	if err != nil {
		retErr := fmt.Errorf("unable to capture all running and pending builds: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}

// swagger:operation PUT /api/v1/admin/build admin AdminUpdateBuild
//
// Update a build in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the build in the database
//     schema:
//       "$ref": "#/definitions/Build"
//   '401':
//     description: Unauthorized to update the build in the database
//     schema:
//       "$ref": "#/definitions/Error
//   '400':
//     description: Unable to update the build in the database - bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the build in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateBuild represents the API handler to
// update any build stored in the database.
func UpdateBuild(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()
	u := user.Retrieve(c)

	logger := logrus.WithFields(logrus.Fields{
		"ip":      util.EscapeValue(c.ClientIP()),
		"path":    util.EscapeValue(c.Request.URL.Path),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	})

	logger.Debug("platform admin: updating build")

	// capture body from API request
	input := new(types.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for build %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"build":    input.GetNumber(),
		"build_id": input.GetID(),
		"repo":     util.EscapeValue(input.GetRepo().GetName()),
		"repo_id":  input.GetRepo().GetID(),
		"org":      util.EscapeValue(input.GetRepo().GetOrg()),
	}).Debug("attempting to update build")

	// send API call to update the build
	b, err := database.FromContext(c).UpdateBuild(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update build %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"build":    b.GetNumber(),
		"build_id": b.GetID(),
		"repo":     b.GetRepo().GetName(),
		"repo_id":  b.GetRepo().GetID(),
		"org":      b.GetRepo().GetOrg(),
	}).Info("updated build")

	c.JSON(http.StatusOK, b)
}
