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

// AllSecrets represents the API handler to
// captures all secrets stored in the database.
func AllSecrets(c *gin.Context) {
	logrus.Info("Reading all secrets")

	// send API call to capture all secrets
	s, err := database.FromContext(c).GetSecretList()
	if err != nil {
		retErr := fmt.Errorf("unable to capture all secrets: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.JSON(http.StatusOK, s)
}
