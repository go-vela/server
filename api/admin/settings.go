// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code
package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

type PlatformSettings_API struct {
	ID     int    `json:"id"`
	FooNum int    `json:"foo_num"`
	BarStr string `json:"bar_str"`
}

// swagger:operation GET /api/v1/admin/settings admin GetSettings
//
// Get the currently configured platform settings.
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved platform settings from the database
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Settings"
//   '500':
//     description: Unable to retrieve platform settings from the database
//     schema:
//       "$ref": "#/definitions/Error"

// GetSettings represents the API handler to
// captures platform settings stored in the database.
func GetSettings(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()

	logrus.Info("Admin: reading platform settings")

	// send API call to capture pending and running builds
	s, err := database.FromContext(c).GetSettings(ctx)
	if err != nil {
		retErr := fmt.Errorf("unable to capture platform settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}

func UpdateSettings(c *gin.Context) {
	// capture middleware values
	// todo: settings.Retrieve
	// s := user.Retrieve(c)

	s, err := database.FromContext(c).GetSettings(c.Request.Context())
	if err != nil {
		retErr := fmt.Errorf("unable to retrieve platform settings from the database: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// maxBuildLimit := c.Value("maxBuildLimit").(int64)
	// defaultRepoEvents := c.Value("defaultRepoEvents").([]string)
	// defaultRepoEventsMask := c.Value("defaultRepoEventsMask").(int64)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.Info("Admin: updating platform settings")

	// capture body from API request
	type ss struct {
		BarStr string `json:"bar_str"`
	}

	input := new(ss)
	// input := new(library.Settings)

	err = c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for platform settings: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// todo: do this right
	updated := "default value"
	if s != nil && input != nil {
		updated = *s + input.BarStr
	}

	// send API call to update the repo
	s, err = database.FromContext(c).UpdateSettings(ctx, &updated)
	if err != nil {
		retErr := fmt.Errorf("unable to update platform settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
