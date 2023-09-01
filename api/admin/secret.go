// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code
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

// swagger:operation PUT /api/v1/admin/secret admin AdminUpdateSecret
//
// Update a secret in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing secret to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Secret"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the secret in the database
//     schema:
//       "$ref": "#/definitions/Secret"
//   '404':
//     description: Unable to update the secret in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '501':
//     description: Unable to update the secret in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateSecret represents the API handler to
// update any secret stored in the database.
func UpdateSecret(c *gin.Context) {
	logrus.Info("Admin: updating secret in database")

	// capture middleware values
	ctx := c.Request.Context()

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the secret
	s, err := database.FromContext(c).UpdateSecret(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}
