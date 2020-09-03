// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//[here] Delete later when confirmed if unneeded.

func finOrgfinder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// oParam := c.Param("org")
		// // c.JSON(http.StatusOK, "hello "+oParam)

		// logrus.Debugf("Reading repo %s", oParam)
		// r, err := database.FromContext(c).GetOrgBuildList(c, 1, 1) //[here] Checks DB for orgs/repo. Sends error if either are not in DB.
		// if err != nil {
		// 	retErr := fmt.Errorf("unable to read repo %s: %v", oParam, err)
		// 	util.HandleError(c, http.StatusNotFound, retErr)
		// 	return
		// }
		// c.JSON(http.StatusOK, r)
		c.JSON(http.StatusOK, "hello")
	}
}

// func Orgfinder(c *gin.Context) {
// 	c.JSON(http.StatusOK, "hello")
// }
