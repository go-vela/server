// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the pipeline in the given context.
func Retrieve(c *gin.Context) *library.Pipeline {
	return FromContext(c)
}

// Establish sets the pipeline in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", c.Param("org"), c.Param("repo"))

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		p := c.Param("pipeline")
		if len(p) == 0 {
			retErr := fmt.Errorf("no pipeline parameter provided")

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"org":      o,
			"pipeline": p,
			"repo":     r.GetName(),
			"user":     u.GetName(),
		}).Debugf("reading pipeline %s/%s", r.GetFullName(), p)

		pipeline, err := database.FromContext(c).GetPipelineForRepo(p, r)
		if err != nil {
			retErr := fmt.Errorf("unable to read pipeline %s/%s: %w", r.GetFullName(), p, err)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		ToContext(c, pipeline)

		c.Next()
	}
}
