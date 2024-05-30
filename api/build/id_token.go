// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/id_token builds GetIDToken
//
// Get a Vela OIDC token associated with a build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
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
//     description: Bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized request
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to find build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to generate id token
//     schema:
//       "$ref": "#/definitions/Error"

// GetIDToken represents the API handler to generate a id token.
func GetIDToken(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	cl := claims.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"org":     b.GetRepo().GetOrg(),
		"repo":    b.GetRepo().GetName(),
		"subject": cl.Subject,
	}).Infof("generating ID token for build %s/%d", b.GetRepo().GetFullName(), b.GetNumber())

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
	if len(c.QueryArray("audience")) > 0 {
		audience := []string{}
		for _, a := range c.QueryArray("audience") {
			if len(a) > 0 {
				audience = append(audience, util.Sanitize(a))
			}
		}
		idmto.Audience = audience
	}

	// mint token
	idt, err := tm.MintIDToken(ctx, idmto, database.FromContext(c))
	if err != nil {
		retErr := fmt.Errorf("unable to generate build token: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &idt})
}
