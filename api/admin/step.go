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

// AllSteps represents the API handler to
// captures all steps stored in the database.
func AllSteps(c *gin.Context) {
	logrus.Info("Reading all steps")

	// send API call to capture all steps
	s, err := database.FromContext(c).GetStepList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all steps: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.JSON(http.StatusOK, s)
}
