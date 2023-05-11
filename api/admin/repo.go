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

// swagger:operation PUT /api/v1/admin/repo admin AdminUpdateRepo
//
// Update a repo in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing repo to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the repo in the database
//     schema:
//       "$ref": "#/definitions/Repo"
//   '404':
//     description: Unable to update the repo in the database
//     schema:
//       "$ref": "#/definitions/Error"
//   '501':
//     description: Unable to update the repo in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateRepo represents the API handler to
// update any repo stored in the database.
func UpdateRepo(c *gin.Context) {
	logrus.Info("Admin: updating repo in database")

	// capture body from API request
	input := new(library.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for repo %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the repo
	err = database.FromContext(c).UpdateRepo(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update repo %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
