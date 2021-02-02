// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the build in the given context.
func Retrieve(c *gin.Context) *library.Build {
	return FromContext(c)
}

// Establish sets the build in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := repo.Retrieve(c)
		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", c.Param("org"), c.Param("repo"))
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		bParam := c.Param("build")
		if len(bParam) == 0 {
			retErr := fmt.Errorf("no build parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		number, err := strconv.Atoi(bParam)
		if err != nil {
			retErr := fmt.Errorf("invalid build parameter provided: %s", bParam)
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		logrus.Debugf("Reading build %s/%d", r.GetFullName(), number)
		b, err := database.FromContext(c).GetBuild(number, r)
		if err != nil {
			retErr := fmt.Errorf("unable to read build %s/%d: %v", r.GetFullName(), number, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		ToContext(c, b)
		c.Next()
	}
}
