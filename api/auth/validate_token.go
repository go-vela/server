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

// swagger:operation POST /authenticate/token/validate authenticate ValidateAuthToken
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

// ValidateAuthToken represents the API handler to
// validate a user access token was created by Vela.
func ValidateAuthToken(c *gin.Context) {
	// capture middleware values
	// ctx := c.Request.Context()

	// attempt to validate access token from source OAuth app
	err := scm.FromContext(c).AuthenticateAccessToken(c.Request)
	if err != nil {
		retErr := fmt.Errorf("unable to authenticate access token: %w", err)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// return a 200 indicating token is valid and created by the server's OAuth app
	c.JSON(http.StatusOK, "access token was created by vela")
}
