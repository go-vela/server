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

// AllUsers represents the API handler to
// captures all users stored in the database.
func AllUsers(c *gin.Context) {
	logrus.Info("Reading all users")

	// send API call to capture all users
	u, err := database.FromContext(c).GetUserList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all users: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.JSON(http.StatusOK, u)
}
