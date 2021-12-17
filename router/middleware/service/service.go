// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the service in the given context.
func Retrieve(c *gin.Context) *library.Service {
	return FromContext(c)
}

// Establish sets the service in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		// capture middleware values
		b := build.Retrieve(c)
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", o, c.Param("repo"))
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		if b == nil {
			retErr := fmt.Errorf("build %s not found for repo %s", c.Param("build"), r.GetFullName())
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		sParam := c.Param("service")
		if len(sParam) == 0 {
			retErr := fmt.Errorf("no service parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		number, err := strconv.Atoi(sParam)
		if err != nil {
			retErr := fmt.Errorf("malformed service parameter provided: %s", sParam)
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"build":   b.GetNumber(),
			"org":     o,
			"service": number,
			"repo":    r.GetName(),
			"user":    u.GetName(),
		}).Debugf("reading service %s/%d/%d", r.GetFullName(), b.GetNumber(), number)

		s, err := database.FromContext(c).GetService(number, b)
		if err != nil {
			// nolint: lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to read service %s/%d/%d: %v", r.GetFullName(), b.GetNumber(), number, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		ToContext(c, s)
		c.Next()
	}
}
