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
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
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
//   '409':
//     description: Conflict (requested id token for build not in running state)
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
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	cl := claims.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  cl.Subject,
	}).Infof("generating ID token for build %s/%d", r.GetFullName(), b.GetNumber())

	// retrieve token manager from context
	tm := c.MustGet("token-manager").(*token.Manager)

	// set mint token options
	idmto := &token.MintTokenOpts{
		BuildNumber:   b.GetNumber(),
		Repo:          r.GetFullName(),
		TokenType:     constants.IDTokenType,
		Commit:        b.GetCommit(),
		TokenDuration: tm.IDTokenDuration,
	}

	// if audience is provided, include that in claims
	if len(c.QueryArray("audience")) > 0 {
		idmto.Audience = c.QueryArray("audience")
	}

	// mint token
	bt, err := tm.MintIDToken(idmto, database.FromContext(c))
	if err != nil {
		retErr := fmt.Errorf("unable to generate build token: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &bt})
}
