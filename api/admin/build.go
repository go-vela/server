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

// AllBuilds represents the API handler to
// captures all builds stored in the database.
func AllBuilds(c *gin.Context) {
	logrus.Info("Admin: reading all builds")

	// send API call to capture all builds
	b, err := database.FromContext(c).GetBuildList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all builds: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}

// UpdateBuild represents the API handler to
// update any build stored in the database.
func UpdateBuild(c *gin.Context) {
	logrus.Info("Admin: updating build in database")

	// capture body from API request
	input := new(library.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for build %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the build
	err = database.FromContext(c).UpdateBuild(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update build %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
