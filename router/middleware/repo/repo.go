// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Retrieve gets the repo in the given context.
func Retrieve(c *gin.Context) *library.Repo {
	return FromContext(c)
}

// Establish sets the repo in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		o := org.Retrieve(c)
		u := user.Retrieve(c)

		rParam := c.Param("repo")
		if len(rParam) == 0 {
			retErr := fmt.Errorf("no repo parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"org":  o,
			"repo": rParam,
			"user": u.GetName(),
		}).Debugf("reading repo %s/%s", o, rParam)

		r, err := database.FromContext(c).GetRepo(o, rParam)
		if err != nil {
			retErr := fmt.Errorf("unable to read repo %s/%s: %v", o, rParam, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		ToContext(c, r)
		c.Next()
	}
}
