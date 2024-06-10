// SPDX-License-Identifier: Apache-2.0

package build

import (
	"errors"
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
// Get a Vela OIDC request token for a build
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
	if len(image) == 0 {
		retErr := errors.New("no step 'image' provided in query parameters")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	request := c.Query("request")
	if len(request) == 0 {
		retErr := errors.New("no 'request' provided in query parameters")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	commands := false

	var err error

	if len(c.Query("commands")) > 0 {
		commands, err = strconv.ParseBool(c.Query("commands"))
		if err != nil {
			retErr := fmt.Errorf("unable to parse 'commands' query parameter as boolean %s: %w", c.Query("commands"), err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	// retrieve token manager from context
	tm := c.MustGet("token-manager").(*token.Manager)

	exp := (time.Duration(b.GetRepo().GetTimeout()) * time.Minute) + tm.BuildTokenBufferDuration

	// set mint token options
	idmto := &token.MintTokenOpts{
		Build:         b,
		Repo:          b.GetRepo().GetFullName(),
		TokenType:     constants.IDRequestTokenType,
		TokenDuration: exp,
		Image:         util.Sanitize(image),
		Request:       util.Sanitize(request),
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
