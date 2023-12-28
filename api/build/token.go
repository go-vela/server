// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/token builds GetBuildToken
//
// Get a build token
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
//     description: Successfully retrieved build token
//     schema:
//       "$ref": "#/definitions/Token"
//   '400':
//     description: Bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '409':
//     description: Conflict (requested build token for build not in pending state)
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to generate build token
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildToken represents the API handler to generate a build token.
func GetBuildToken(c *gin.Context) {
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
	}).Infof("generating build token for build %s/%d", r.GetFullName(), b.GetNumber())

	// if build is not in a pending state, then a build token should not be needed - conflict
	if !strings.EqualFold(b.GetStatus(), constants.StatusPending) {
		retErr := fmt.Errorf("unable to mint build token: build is not in pending state")
		util.HandleError(c, http.StatusConflict, retErr)

		return
	}

	// retrieve token manager from context
	tm := c.MustGet("token-manager").(*token.Manager)

	// set expiration to repo timeout plus configurable buffer
	exp := (time.Duration(r.GetTimeout()) * time.Minute) + tm.BuildTokenBufferDuration

	// set mint token options
	bmto := &token.MintTokenOpts{
		Hostname:      cl.Subject,
		BuildID:       b.GetID(),
		Repo:          r.GetFullName(),
		TokenType:     constants.WorkerBuildTokenType,
		TokenDuration: exp,
	}

	// mint token
	bt, err := tm.MintToken(bmto)
	if err != nil {
		retErr := fmt.Errorf("unable to generate build token: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &bt})
}
