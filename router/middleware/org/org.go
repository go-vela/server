// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package org

import (
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Retrieve gets the repo in the given context
// func Retrieve(c *gin.Context) *library.Repo { //[here] if string/point doesn't work, try creating an org struct
// 	return FromContext(c)
// }

// Establish sets the org in the given context
func Establish() gin.HandlerFunc { //[here] Note: Without this function, the data will not parse and API returns an empty set.
	return func(c *gin.Context) {
		oParam := c.Param("org")
		if len(oParam) == 0 { //[here] Checks if org is present. Sends error if not. Note: May not be needed.
			retErr := fmt.Errorf("no org parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		//[here] This is technically not used.
		logrus.Debugf("Reading org %s/%s", oParam)
		o, err := database.FromContext(c).GetRepoOrg(oParam) //[here] Checks DB for orgs. Sends error if not in DB.
		if err != nil {
			retErr := fmt.Errorf("unable to read org %s: %v", oParam, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		ToContext(c, o) //[here] 100% needs to be changes. Uses repo struct.
		c.Next()
	}
}
