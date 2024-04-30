// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/id_request_token builds GetIDRequestToken
//
// Get a Vela OIDC request token associated with a build
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
//   name: image
//   description: Add image to token claims
//   type: string
// - in: query
//   name: request
//   description: Add request input to token claims
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved ID Request token
//     schema:
//       "$ref": "#/definitions/Token"
//   '400':
//     description: Bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to generate ID request token
//     schema:
//       "$ref": "#/definitions/Error"

// GetIDRequestToken represents the API handler to generate and return an ID request token.
func GetIDRequestToken(c *gin.Context) {
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
	}).Infof("generating ID request token for build %s/%d", r.GetFullName(), b.GetNumber())

	image := c.Query("image")
	request := c.Query("request")

	// retrieve token manager from context
	tm := c.MustGet("token-manager").(*token.Manager)

	exp := (time.Duration(r.GetTimeout()) * time.Minute) + tm.BuildTokenBufferDuration

	// set mint token options
	idmto := &token.MintTokenOpts{
		Build:         b,
		Repo:          r.GetFullName(),
		TokenType:     constants.IDRequestTokenType,
		Commit:        b.GetCommit(),
		TokenDuration: exp,
		Image:         image,
		Request:       request,
	}

	// mint token
	bt, err := tm.MintToken(idmto)
	if err != nil {
		retErr := fmt.Errorf("unable to generate ID request token: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &bt})
}
