// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/search/builds/{id} builds GetBuildByID
//
// Get a single build by its id in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: id
//   description: build id
//   required: true
//   type: number
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved build
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to retrieve the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the build
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildByID represents the API handler to capture a
// build by its id from the configured backend.
func GetBuildByID(c *gin.Context) {
	// Variables that will hold the library types of the build and repo
	var (
		b *library.Build
		r *library.Repo
	)

	// Capture user from middleware
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// Parse build ID from path
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		retErr := fmt.Errorf("unable to parse build id: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": id,
		"user":  u.GetName(),
	}).Infof("reading build %d", id)

	// Get build from database
	b, err = database.FromContext(c).GetBuild(ctx, id)
	if err != nil {
		retErr := fmt.Errorf("unable to get build: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Get repo from database using repo ID field from build
	r, err = database.FromContext(c).GetRepo(ctx, b.GetRepoID())
	if err != nil {
		retErr := fmt.Errorf("unable to get repo: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Capture user access from SCM. We do this in order to ensure user has access and is not
	// just retrieving any build using a random id number.
	perm, err := scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
	}

	// Ensure that user has at least read access to repo to return the build
	if perm == "none" && !u.GetAdmin() {
		retErr := fmt.Errorf("unable to retrieve build %d: user does not have read access to repo", id)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}
