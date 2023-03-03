// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the initstep in the given context.
func Retrieve(c *gin.Context) *library.InitStep {
	return FromContext(c)
}

// Establish sets the initstep in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		// capture middleware values
		b := build.Retrieve(c)
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", util.PathParameter(c, "org"), util.PathParameter(c, "repo"))

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		if b == nil {
			retErr := fmt.Errorf("build %s not found for repo %s", util.PathParameter(c, "build"), r.GetFullName())
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		initStepParam := util.PathParameter(c, "initstep")
		if len(initStepParam) == 0 {
			retErr := fmt.Errorf("no initstep parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		number, err := strconv.Atoi(initStepParam)
		if err != nil {
			retErr := fmt.Errorf("malformed initstep parameter provided: %s", initStepParam)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"org":      o,
			"repo":     r.GetName(),
			"build":    b.GetNumber(),
			"initstep": number,
			"user":     u.GetName(),
		}).Debugf("reading initstep %s/%d/%d", r.GetFullName(), b.GetNumber(), number)

		s, err := database.FromContext(c).GetInitStepForBuild(b, number)
		if err != nil {
			retErr := fmt.Errorf("unable to read initstep %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), number, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		ToContext(c, s)
		c.Next()
	}
}
