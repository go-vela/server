// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package admin

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/admin/secrets admin AdminAllSecrets
//
// Get all of the secrets in the database
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved all secrets from the database
//     type: json
//     schema:
//       "$ref": "#/definitions/Secret"
//   '500':
//     description: Unable to retrieve all secrets from the database
//     schema:
//       type: string

// AllSecrets represents the API handler to
// captures all secrets stored in the database.
func AllSecrets(c *gin.Context) {
	logrus.Info("Admin: reading all secrets")

	// send API call to capture all secrets
	s, err := database.FromContext(c).GetSecretList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all secrets: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}

// swagger:operation PUT /api/v1/admin/secret admin AdminUpdateSecret
//
// Update a secret in the database
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing secret to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Secret"
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully updated the secret in the database
//     type: json
//     schema:
//       "$ref": "#/definitions/Secret"
//   '404':
//     description: Unable to update the secret in the database
//     schema:
//       type: string
//   '501':
//     description: Unable to update the secret in the database
//     schema:
//       type: string

// UpdateSecret represents the API handler to
// update any secret stored in the database.
func UpdateSecret(c *gin.Context) {
	logrus.Info("Admin: updating secret in database")

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the secret
	err = database.FromContext(c).UpdateSecret(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
