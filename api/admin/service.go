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

// AllServices represents the API handler to
// captures all services stored in the database.
func AllServices(c *gin.Context) {
	logrus.Info("Admin: reading all services")

	// send API call to capture all services
	s, err := database.FromContext(c).GetServiceList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all services: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s)
}

// UpdateService represents the API handler to
// update any service stored in the database.
func UpdateService(c *gin.Context) {
	logrus.Info("Admin: updating service in database")

	// capture body from API request
	input := new(library.Service)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the service
	err = database.FromContext(c).UpdateService(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update service %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
