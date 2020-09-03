// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package org

import (
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Retrieve gets the repo in the given context
func Retrieve(c *gin.Context) *library.Repo {
	return FromContext(c)
}

// Establish sets the repo in the given context
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		oParam := c.Param("org")
		if len(oParam) == 0 { //[here] Checks if org is present. Sends error if not. Note: May not be needed.
			retErr := fmt.Errorf("no org parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		logrus.Debugf("Reading repo %s/%s", oParam)
		r, err := database.FromContext(c).GetRepoOrg(oParam) //[here] Checks DB for orgs. Sends error if not in DB.
		if err != nil {
			retErr := fmt.Errorf("unable to read org %s: %v", oParam, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		ToContext(c, r)
		c.Next()
	}
}
