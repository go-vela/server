// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/image"
	"github.com/go-vela/server/queue"
	cliMiddleware "github.com/go-vela/server/router/middleware/cli"
	sMiddleware "github.com/go-vela/server/router/middleware/settings"
	uMiddleware "github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/admin/settings admin GetSettings
//
// Get platform settings
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved settings
//     type: json
//     schema:
//       "$ref": "#/definitions/Platform"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"

// GetSettings represents the API handler to get platform settings.
func GetSettings(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)

	l.Debug("platform admin: reading platform settings")

	// capture middleware values
	s := sMiddleware.FromContext(c)

	// check captured value because we aren't retrieving settings from the database
	// instead we are retrieving the auto-refreshed middleware value
	if s == nil {
		retErr := fmt.Errorf("platform settings not found")

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}

// swagger:operation PUT /api/v1/admin/settings admin UpdateSettings
//
// Update platform settings
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The settings object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Platform"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated platform settings
//     type: json
//     schema:
//       "$ref": "#/definitions/Platform"
//   '400':
//     description: Invalid request payload
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

// UpdateSettings represents the API handler to update the
// platform settings singleton.
//
//nolint:funlen // nil checks throughout handler make this function long
func UpdateSettings(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	s := sMiddleware.FromContext(c)
	u := uMiddleware.FromContext(c)
	ctx := c.Request.Context()

	l.Debug("platform admin: updating platform settings")

	// check captured value because we aren't retrieving settings from the database
	// instead we are retrieving the auto-refreshed middleware value
	if s == nil {
		retErr := fmt.Errorf("platform settings not found")

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// duplicate settings to not alter the shared pointer
	_s := new(settings.Platform)
	_s.FromSettings(s)

	// ensure we update the singleton record
	_s.SetID(1)

	// capture body from API request
	input := new(settings.Platform)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for platform settings: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if input.Compiler != nil {
		if input.CloneImage != nil {
			// validate clone image
			cloneImage := *input.CloneImage

			_, err = image.ParseWithError(cloneImage)
			if err != nil {
				retErr := fmt.Errorf("invalid clone image %s for platform settings: %w", cloneImage, err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			_s.SetCloneImage(cloneImage)

			l.Infof("platform admin: updating clone image to %s", cloneImage)
		}

		if input.TemplateDepth != nil {
			_s.SetTemplateDepth(*input.TemplateDepth)

			l.Infof("platform admin: updating template depth to %d", *input.TemplateDepth)
		}

		if input.StarlarkExecLimit != nil {
			_s.SetStarlarkExecLimit(*input.StarlarkExecLimit)

			l.Infof("platform admin: updating starlark exec limit to %d", *input.StarlarkExecLimit)
		}
	}

	if input.Queue != nil {
		if input.Queue.Routes != nil {
			_s.SetRoutes(input.GetRoutes())
		}

		l.Infof("platform admin: updating queue routes to: %s", input.GetRoutes())
	}

	if input.SCM != nil {
		if input.SCM.RepoRoleMap != nil {
			err = util.ValidateRoleMap(input.GetRepoRoleMap(), "repo")
			if err != nil {
				retErr := fmt.Errorf("invalid repo role map for platform settings: %w", err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			_s.SetRepoRoleMap(input.GetRepoRoleMap())

			l.Infof("platform admin: updating repo role map to: %s", input.GetRepoRoleMap())
		}

		if input.SCM.OrgRoleMap != nil {
			err = util.ValidateRoleMap(input.GetOrgRoleMap(), "org")
			if err != nil {
				retErr := fmt.Errorf("invalid org role map for platform settings: %w", err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			_s.SetOrgRoleMap(input.GetOrgRoleMap())

			l.Infof("platform admin: updating org role map to: %s", input.GetOrgRoleMap())
		}

		if input.SCM.TeamRoleMap != nil {
			err = util.ValidateRoleMap(input.GetTeamRoleMap(), "team")
			if err != nil {
				retErr := fmt.Errorf("invalid team role map for platform settings: %w", err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			_s.SetTeamRoleMap(input.GetTeamRoleMap())

			l.Infof("platform admin: updating team role map to: %s", input.GetTeamRoleMap())
		}
	}

	if input.RepoAllowlist != nil {
		_s.SetRepoAllowlist(input.GetRepoAllowlist())

		l.Infof("platform admin: updating repo allowlist to: %s", input.GetRepoAllowlist())
	}

	if input.ScheduleAllowlist != nil {
		_s.SetScheduleAllowlist(input.GetScheduleAllowlist())

		l.Infof("platform admin: updating schedule allowlist to: %s", input.GetScheduleAllowlist())
	}

	if input.MaxDashboardRepos != nil {
		_s.SetMaxDashboardRepos(input.GetMaxDashboardRepos())

		l.Infof("platform admin: updating max dashboard repos to: %d", input.GetMaxDashboardRepos())
	}

	_s.SetUpdatedBy(u.GetName())

	// send API call to update the settings
	_s, err = database.FromContext(c).UpdateSettings(ctx, _s)
	if err != nil {
		retErr := fmt.Errorf("unable to update platform settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, _s)
}

// swagger:operation DELETE /api/v1/admin/settings admin RestoreSettings
//
// Restore platform settings to the environment defaults
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully restored default platform settings
//     type: json
//     schema:
//       "$ref": "#/definitions/Platform"
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

// RestoreSettings represents the API handler to
// restore platform settings to the environment defaults.
func RestoreSettings(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)

	l.Debug("platform admin: restoring platform settings")

	// capture middleware values
	ctx := c.Request.Context()
	cliCmd := cliMiddleware.FromContext(c)
	s := sMiddleware.FromContext(c)
	u := uMiddleware.FromContext(c)

	// check captured value because we aren't retrieving settings from the database
	// instead we are retrieving the auto-refreshed middleware value
	if s == nil {
		retErr := fmt.Errorf("platform settings not found")

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	compiler, err := native.FromCLICommand(ctx, cliCmd)
	if err != nil {
		retErr := fmt.Errorf("unable to restore platform settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	queue, err := queue.FromCLICommand(ctx, cliCmd)
	if err != nil {
		retErr := fmt.Errorf("unable to restore platform settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// initialize a new settings record
	_s := settings.FromCLICommand(cliCmd)

	_s.SetUpdatedAt(time.Now().UTC().Unix())
	_s.SetUpdatedBy(u.GetName())

	// read in defaults supplied from the cli runtime
	compilerSettings := compiler.GetSettings()
	_s.SetCompiler(compilerSettings)

	queueSettings := queue.GetSettings()
	_s.SetQueue(queueSettings)

	// send API call to update the settings
	s, err = database.FromContext(c).UpdateSettings(ctx, _s)
	if err != nil {
		retErr := fmt.Errorf("unable to update (restore) platform settings: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
