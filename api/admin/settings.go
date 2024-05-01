// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/database"
	sMiddleware "github.com/go-vela/server/router/middleware/settings"
	uMiddleware "github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/admin/settings admin GetSettings
//
// Get the currently configured settings.
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved settings from the database
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Platform"
//   '500':
//     description: Unable to retrieve settings from the database
//     schema:
//       "$ref": "#/definitions/Error"

// GetSettings represents the API handler to
// captures settings stored in the database.
func GetSettings(c *gin.Context) {
	// capture middleware values
	s := sMiddleware.FromContext(c)

	logrus.Info("Admin: reading settings")

	c.JSON(http.StatusOK, s)
}

// swagger:operation PUT /api/v1/admin/settings admin UpdateSettings
//
// Update the settings singleton in the database.
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing settings to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Platform"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the settings in the database
//     schema:
//       "$ref": "#/definitions/Platform"
//   '404':
//     description: Unable to update the settings in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '501':
//     description: Unable to update the settings in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateSettings represents the API handler to
// update the settings singleton stored in the database.
func UpdateSettings(c *gin.Context) {
	// capture middleware values
	s := sMiddleware.FromContext(c)
	u := uMiddleware.FromContext(c)
	ctx := c.Request.Context()

	logrus.Info("Admin: updating settings")

	// capture body from API request
	input := new(settings.Platform)

	input.SetUpdatedAt(time.Now().UTC().Unix())
	input.SetUpdatedBy(u.GetName())

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for settings: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if input.Compiler != nil {
		if input.CloneImage != nil {
			s.CloneImage = input.CloneImage
		}

		if input.TemplateDepth != nil {
			s.TemplateDepth = input.TemplateDepth
		}

		if input.StarlarkExecLimit != nil {
			s.StarlarkExecLimit = input.StarlarkExecLimit
		}
	}

	if input.Queue != nil {
		if input.Queue.Routes != nil {
			s.SetRoutes(input.GetRoutes())
		}
	}

	if input.RepoAllowlist != nil {
		s.SetRepoAllowlist(input.GetRepoAllowlist())
	}

	if input.ScheduleAllowlist != nil {
		s.SetScheduleAllowlist(input.GetScheduleAllowlist())
	}

	// send API call to update the settings
	s, err = database.FromContext(c).UpdateSettings(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to update settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
