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

// AllSteps represents the API handler to
// captures all steps stored in the database.
func AllSteps(c *gin.Context) {
	logrus.Info("Admin: reading all steps")

	// send API call to capture all steps
	s, err := database.FromContext(c).GetStepList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all steps: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}

// UpdateStep represents the API handler to
// update any step stored in the database.
func UpdateStep(c *gin.Context) {
	logrus.Info("Admin: updating step in database")

	// capture body from API request
	input := new(library.Step)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for step %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the step
	err = database.FromContext(c).UpdateStep(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update step %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
