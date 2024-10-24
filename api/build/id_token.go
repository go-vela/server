// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/id_token builds GetIDToken
//
// Get a Vela OIDC token for a build
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
//   name: audience
//   description: Add audience to token claims
//   type: array
//   items:
//     type: string
//   collectionFormat: multi
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved ID token
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

// GetIDToken represents the API handler to generate a id token.
func GetIDToken(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	cl := claims.Retrieve(c)
	ctx := c.Request.Context()

	l.Infof("generating ID token for build %s/%d", b.GetRepo().GetFullName(), b.GetNumber())

	// retrieve token manager from context
	tm := c.MustGet("token-manager").(*token.Manager)

	// set mint token options
	idmto := &token.MintTokenOpts{
		Build:         b,
		Repo:          b.GetRepo().GetFullName(),
		TokenType:     constants.IDTokenType,
		TokenDuration: tm.IDTokenDuration,
		Image:         cl.Image,
		Request:       cl.Request,
		Commands:      cl.Commands,
	}

	// if audience is provided, include that in claims
	audience := []string{}

	if len(c.QueryArray("audience")) > 0 {
		for _, a := range c.QueryArray("audience") {
			if len(a) > 0 {
				audience = append(audience, util.Sanitize(a))
			}
		}
	}

	if len(audience) == 0 {
		retErr := fmt.Errorf("unable to generate ID token: %s", "no audience provided")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	idmto.Audience = audience

	// mint token
	idt, err := tm.MintIDToken(ctx, idmto, database.FromContext(c))
	if err != nil {
		retErr := fmt.Errorf("unable to generate ID token: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, types.Token{Token: &idt})
}
