// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
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
//     description: Unable to generate ID request token
//     schema:
//       "$ref": "#/definitions/Error"

// GetIDRequestToken represents the API handler to generate and return an ID request token.
func GetIDRequestToken(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	cl := claims.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   b.GetRepo().GetOrg(),
		"repo":  b.GetRepo().GetName(),
		"user":  cl.Subject,
	}).Infof("generating ID request token for build %s/%d", b.GetRepo().GetFullName(), b.GetNumber())

	image := c.Query("image")
	request := c.Query("request")
	commands, _ := strconv.ParseBool(c.Query("commands"))

	// retrieve token manager from context
	tm := c.MustGet("token-manager").(*token.Manager)

	exp := (time.Duration(b.GetRepo().GetTimeout()) * time.Minute) + tm.BuildTokenBufferDuration

	// set mint token options
	idmto := &token.MintTokenOpts{
		Build:         b,
		Repo:          b.GetRepo().GetFullName(),
		TokenType:     constants.IDRequestTokenType,
		TokenDuration: exp,
		Image:         image,
		Request:       request,
		Commands:      commands,
	}

	// mint token
	idrt, err := tm.MintToken(idmto)
	if err != nil {
		retErr := fmt.Errorf("unable to generate ID request token: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &idrt})
}
