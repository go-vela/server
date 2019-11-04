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

// AllBuilds represents the API handler to
// captures all builds stored in the database.
func AllBuilds(c *gin.Context) {
	logrus.Info("Reading all builds")

	// send API call to capture all builds
	b, err := database.FromContext(c).GetBuildList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all builds: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.JSON(http.StatusOK, b)
}
