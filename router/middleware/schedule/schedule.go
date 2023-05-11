// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"fmt"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/util"
)

// Retrieve gets the schedule in the given context.
func Retrieve(c *gin.Context) *types.Schedule {
	return FromContext(c)
}

// Establish sets the schedule in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		sParam := util.PathParameter(c, "schedule")
		if len(sParam) == 0 {
			retErr := fmt.Errorf("no schedule parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"org":  r.GetOrg(),
			"repo": r.GetName(),
			"user": u.GetName(),
		}).Debugf("reading schedule %s for repo %s", sParam, r.GetFullName())

		s, err := database.FromContext(c).GetScheduleForRepo(r, sParam)
		if err != nil {
			retErr := fmt.Errorf("unable to read schedule %s for repo %s: %w", sParam, r.GetFullName(), err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		ToContext(c, s)
		c.Next()
	}
}
