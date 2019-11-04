// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the step in the given context
func Retrieve(c *gin.Context) *library.Step {
	return FromContext(c)
}

// Establish sets the step in the given context
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := repo.Retrieve(c)
		if r == nil {
			retErr := fmt.Errorf("Repo %s/%s not found", c.Param("org"), c.Param("repo"))
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		b := build.Retrieve(c)
		if b == nil {
			retErr := fmt.Errorf("Build %s not found for repo %s/%s", c.Param("build"), c.Param("org"), c.Param("repo"))
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		sParam := c.Param("step")
		if len(sParam) == 0 {
			retErr := fmt.Errorf("No step parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		number, err := strconv.Atoi(sParam)
		if err != nil {
			retErr := fmt.Errorf("Malformed step parameter provided: %s", sParam)
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}

		logrus.Debugf("Reading step %s/%d/%d", r.GetFullName(), b.Number, number)
		s, err := database.FromContext(c).GetStep(number, b)
		if err != nil {
			retErr := fmt.Errorf("Error while reading step %s/%d/%d: %v", r.GetFullName(), b.Number, number, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		ToContext(c, s)
		c.Next()
	}
}
