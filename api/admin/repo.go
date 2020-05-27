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

// swagger:operation GET /api/v1/admin/repos admin AllRepos
//
// Get all of the repos in the database
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved all repos from the database
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"
//   '500':
//     description: Unable to retrieve all repos from the database
//     schema:
//       type: string

// AllRepos represents the API handler to
// captures all repos stored in the database.
func AllRepos(c *gin.Context) {
	logrus.Info("Admin: reading all repos")

	// send API call to capture all repos
	r, err := database.FromContext(c).GetRepoList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all repos: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, r)
}

// swagger:operation PUT /api/v1/admin/repo admin UpdateRepo
//
// Update a repo in the database
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing repo to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully updated the repo in the database
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"
//   '404':
//     description: unable to update the repo in the database
//     schema:
//       type: string
//   '501':
//     description: Unable to update the repo in the database
//     schema:
//       type: string

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
