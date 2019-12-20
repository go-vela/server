// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package admin

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AllRepos represents the API handler to
// captures all repos stored in the database.
func AllRepos(c *gin.Context) {
	logrus.Info("Reading all repos")

	// send API call to capture all repos
	r, err := database.FromContext(c).GetRepoList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all repos: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, r)
}
