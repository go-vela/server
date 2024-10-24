// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /authenticate/token authenticate PostAuthToken
//
// Authenticate to Vela via personal access token
//
// ---
// produces:
// - application/json
// parameters:
// - in: header
//   name: Token
//   type: string
//   required: true
//   description: >
//     scopes: repo, repo:status, user:email, read:user, and read:org
// responses:
//   '200':
//     description: Successfully authenticated
//     schema:
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Service unavailable
//     schema:
//       "$ref": "#/definitions/Error"

// PostAuthToken represents the API handler to process
// a user logging in using PAT to Vela from the API.
func PostAuthToken(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	// attempt to get user from source
	u, err := scm.FromContext(c).AuthenticateToken(ctx, c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to authenticate user: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	l.Infof("SCM user %s authenticated using PAT", u.GetName())

	// check if the user exists
	u, err = database.FromContext(c).GetUserForName(ctx, u.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to authenticate: user %s not found", u.GetName())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"user":    u.GetName(),
		"user_id": u.GetID(),
	}).Info("user successfully authenticated via SCM PAT")

	// We don't need refresh token for this scenario
	// We only need access token and are configured based on the config defined
	tm := c.MustGet("token-manager").(*token.Manager)

	// mint token options for access token
	amto := &token.MintTokenOpts{
		User:          u,
		TokenType:     constants.UserAccessTokenType,
		TokenDuration: tm.UserAccessTokenDuration,
	}

	at, err := tm.MintToken(amto)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)
	}

	// return jwt access token
	c.JSON(http.StatusOK, types.Token{Token: &at})
}
