// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/settings"
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
	s := settings.FromContext(c)

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

// todo: swagger and comments
func UpdateSettings(c *gin.Context) {
	// capture middleware values
	s := settings.FromContext(c)
	ctx := c.Request.Context()

	logrus.Info("Admin: updating settings")

	// capture body from API request
	input := new(api.Settings)

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

	if input.QueueRoutes != nil {
		// update routes if set
		s.SetQueueRoutes(input.GetQueueRoutes())
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
