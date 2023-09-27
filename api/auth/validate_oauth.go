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

// swagger:operation GET /validate-oauth authenticate ValidateOAuthToken
//
// Validate that a user oauth token was created by Vela
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
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unable to validate
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Service unavailable
//     schema:
//       "$ref": "#/definitions/Error"

// ValidateOAuthToken represents the API handler to
// validate that a user oauth token was created by Vela.
func ValidateOAuthToken(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()

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
