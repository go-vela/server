// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /authenticate/token/validate authenticate ValidateOAuthToken
//
// Validate that a user access token was created by Vela
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
//     description: Successfully authenticated
//     schema:
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unable to authenticate
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Service unavailable
//     schema:
//       "$ref": "#/definitions/Error"

// ValidateOAuthToken represents the API handler to
// validate a user oauth token was created by Vela.
func ValidateOAuthToken(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()

	token := c.Request.Header.Get("Token")
	if len(token) == 0 {
		retErr := fmt.Errorf("unable to validate access token: no token provided in header")

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// attempt to validate access token from source OAuth app
	err := scm.FromContext(c).ValidateOAuthToken(ctx, token)
	if err != nil {
		retErr := fmt.Errorf("unable to validate oauth token: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// return a 200 indicating token is valid and created by the server's OAuth app
	c.JSON(http.StatusOK, "oauth token was created by vela")
}
