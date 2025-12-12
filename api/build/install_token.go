// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/router/middleware/auth"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/install_token builds GetInstallToken
//
// Get a Vela GitHub App install token
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
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: query
//   name: image
//   description: Add image to token claims
//   type: string
// - in: query
//   name: request
//   description: Add request input to token claims
//   type: string
// - in: query
//   name: commands
//   description: Add commands input to token claims
//   type: boolean
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved ID Request token
//     schema:
//       "$ref": "#/definitions/Token"
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

// GetInstallToken represents the API handler to generate and return an install token.
func GetInstallToken(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	ctx := c.Request.Context()
	tknCache := cache.FromContext(c)

	l.Debugf("generating install token for build %s/%d", b.GetRepo().GetFullName(), b.GetNumber())

	// build must be running to refresh install token
	if b.GetStatus() != constants.StatusRunning {
		retErr := fmt.Errorf("unable to generate install token for build not in %s status", constants.StatusRunning)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	installToken := auth.RetrieveTokenHeader(c.Request)
	if installToken == "" {
		retErr := fmt.Errorf("unable to retrieve installation token from request header")

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// fetch input token from cache
	cachedToken, err := tknCache.GetInstallToken(ctx, installToken)
	if err != nil {
		retErr := fmt.Errorf("unable to get installation token from cache: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// if cached token is still valid for 5+ minutes, return it
	if cachedToken.Expiration > time.Now().Add(5*time.Minute).Unix() {
		l.Debugf("returning cached install token for build %s/%d", b.GetRepo().GetFullName(), b.GetNumber())

		resp := new(types.Token)
		resp.SetToken(cachedToken.Token)
		resp.SetExpiration(cachedToken.Expiration)

		c.JSON(http.StatusOK, resp)

		return
	}

	// mint new token with same permissions and repositories
	newToken, _, err := scm.FromContext(c).NewAppInstallationToken(ctx, b.GetRepo(), cachedToken.Repositories, cachedToken.Permissions)
	if err != nil {
		retErr := fmt.Errorf("unable to generate new installation token: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// evict old token from cache
	err = tknCache.EvictInstallToken(ctx, cachedToken.Token)
	if err != nil {
		l.Warnf("unable to evict installation token from cache: %v", err)
	}

	// store new token in cache
	err = tknCache.StoreInstallToken(ctx, newToken, b.GetRepo())
	if err != nil {
		retErr := fmt.Errorf("unable to store installation token in cache: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	resp := new(types.Token)
	resp.SetToken(newToken.Token)
	resp.SetExpiration(newToken.Expiration)

	c.JSON(http.StatusOK, resp)
}
