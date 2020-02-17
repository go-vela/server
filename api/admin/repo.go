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
