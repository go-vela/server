// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /validate-oauth authenticate ValidateOAuthToken
//
// Validate that a user OAuth token was created by Vela
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
//     OAuth integration user access token
// responses:
//   '200':
//     description: Successfully validated
//     schema:
//       type: string
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// ValidateOAuthToken represents the API handler to
// validate that a user OAuth token was created by Vela.
func ValidateOAuthToken(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Info("validating oauth token")

	token := c.Request.Header.Get("Token")
	if len(token) == 0 {
		retErr := fmt.Errorf("unable to validate oauth token: no token provided in header")

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// attempt to validate access token from source OAuth app
	ok, err := scm.FromContext(c).ValidateOAuthToken(ctx, token)
	if err != nil {
		retErr := fmt.Errorf("unable to validate oauth token: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	if !ok {
		retErr := fmt.Errorf("oauth token was not created by vela")

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// return a 200 indicating token is valid and created by the server's OAuth app
	c.JSON(http.StatusOK, "oauth token was created by vela")
}
