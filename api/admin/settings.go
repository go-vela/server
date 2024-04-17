// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	settings "github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/database"
	sMiddleware "github.com/go-vela/server/router/middleware/settings"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
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
//         "$ref": "#/definitions/Settings"
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

	output := strings.ToLower(c.Query("output"))

	switch output {
	case "env":
		exported := s.ToEnv()
		c.String(http.StatusOK, exported)
	case "yaml":
		exported := s.ToYAML()
		c.String(http.StatusOK, exported)
	case "json":
		fallthrough
	default:
		c.JSON(http.StatusOK, s)
	}
}

// swagger:operation PUT /api/v1/admin/settings admin AdminUpdateSettings
//
// Update the settings singleton in the database
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
//     "$ref": "#/definitions/Settings"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the settings in the database
//     schema:
//       "$ref": "#/definitions/Settings"
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
	ctx := c.Request.Context()

	logrus.Info("Admin: updating settings")

	// capture body from API request
	input := new(settings.Platform)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for settings: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if input.CloneImage != nil {
		s.CloneImage = input.CloneImage
	}

	if input.TemplateDepth != nil {
		s.TemplateDepth = input.TemplateDepth
	}

	if input.StarlarkExecLimit != nil {
		s.StarlarkExecLimit = input.StarlarkExecLimit
	}

	if input.Routes != nil {
		// update routes if set
		s.SetRoutes(input.GetRoutes())
	}

	if input.RepoAllowlist != nil {
		// update allowlist if set
		s.SetRepoAllowlist(input.GetRepoAllowlist())
	}

	// send API call to update the repo
	s, err = database.FromContext(c).UpdateSettings(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to update settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
