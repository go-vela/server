// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/git_token build PostGitToken
//
// Generate a scoped git token for a build
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
// - in: body
//   name: body
//   description: Token request
//   required: true
//   schema:
//     "$ref": "#/definitions/TokenRequest"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully generated install token
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
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// PostGitToken represents the API handler to generate and return a git token.
func PostGitToken(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	ctx := c.Request.Context()

	if b.GetRepo().GetInstallID() == 0 {
		retErr := fmt.Errorf("repository does not have an installation ID, cannot generate git token")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	l.Debugf("generating install token for build %s/%d", b.GetRepo().GetFullName(), b.GetNumber())

	// build must be running to generate install token
	if b.GetStatus() != constants.StatusRunning {
		retErr := fmt.Errorf("unable to generate install token for build not in %s status", constants.StatusRunning)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture body from API request
	input := new(types.TokenRequest)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for token request for build %s/%d: %w", b.GetRepo().GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// mint new token
	newToken, err := scm.FromContext(c).GetNetrcPassword(ctx, b, input.Repositories, input.Permissions)
	if err != nil {
		retErr := fmt.Errorf("unable to generate new installation token: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	err = cache.FromContext(c).StoreInstallToken(ctx, newToken, b.GetID(), b.GetRepo().GetTimeout())
	if err != nil {
		retErr := fmt.Errorf("unable to store installation token in cache for build %d: %w", b.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	resp := new(types.Token)
	resp.SetToken(newToken)

	c.JSON(http.StatusOK, resp)
}
