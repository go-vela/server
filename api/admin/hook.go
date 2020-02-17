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

// AllHooks represents the API handler to
// captures all hooks stored in the database.
func AllHooks(c *gin.Context) {
	logrus.Info("Admin: reading all hooks")

	// send API call to capture all hooks
	r, err := database.FromContext(c).GetHookList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all hooks: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, r)
}

// UpdateHook represents the API handler to
// update any hook stored in the database.
func UpdateHook(c *gin.Context) {
	logrus.Info("Admin: updating hook in database")

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for hook %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to update the hook
	err = database.FromContext(c).UpdateHook(input)
	if err != nil {
		retErr := fmt.Errorf("unable to update hook %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, input)
}
