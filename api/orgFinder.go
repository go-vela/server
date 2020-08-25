// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

func Orgfinder() gin.HandlerFunc {
	return func(c *gin.Context) {
		oParam := c.Param("org")
		// c.JSON(http.StatusOK, "hello "+oParam)
		logrus.Debugf("Reading repo %s", oParam)
		r, err := database.FromContext(c).GetRepoOrg(oParam) //[here] Checks DB for orgs/repo. Sends error if either are not in DB.
		if err != nil {
			retErr := fmt.Errorf("unable to read repo %s: %v", oParam, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}
		c.JSON(http.StatusOK, r)

	}
}

// func Orgfinder(c *gin.Context) {
// 	c.JSON(http.StatusOK, "hello")
// }
