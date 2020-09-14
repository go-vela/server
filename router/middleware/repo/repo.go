// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

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
		if len(oParam) == 0 {
			retErr := fmt.Errorf("no org parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		rParam := c.Param("repo")
		if len(rParam) == 0 {
			retErr := fmt.Errorf("no repo parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		logrus.Debugf("Reading repo %s/%s", oParam, rParam)
		r, err := database.FromContext(c).GetRepo(oParam, rParam)
		if err != nil {
			retErr := fmt.Errorf("unable to read repo %s/%s: %v", oParam, rParam, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		ToContext(c, r)
		c.Next()
	}
}

// EstablishLite used to check if org param is used only.
func EstablishLite() gin.HandlerFunc { //Note: Func same subject to change and/or moved elsewhere.
	return func(c *gin.Context) {
		oParam := c.Param("org")
		if len(oParam) == 0 {
			retErr := fmt.Errorf("no org parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}
		c.Next()
	}
}
