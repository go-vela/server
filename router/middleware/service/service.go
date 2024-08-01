// SPDX-License-Identifier: Apache-2.0

package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// Retrieve gets the service in the given context.
func Retrieve(c *gin.Context) *library.Service {
	return FromContext(c)
}

// Establish sets the service in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		// capture middleware values
		l := c.MustGet("logger").(*logrus.Entry)
		b := build.Retrieve(c)
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		ctx := c.Request.Context()

		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", o, util.PathParameter(c, "repo"))
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		if b == nil {
			retErr := fmt.Errorf("build %s not found for repo %s", util.PathParameter(c, "build"), r.GetFullName())
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		sParam := util.PathParameter(c, "service")
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

		l.Debugf("reading service %d", number)

		s, err := database.FromContext(c).GetServiceForBuild(ctx, b, number)
		if err != nil {
			retErr := fmt.Errorf("unable to read service %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), number, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"service":    s.GetName(),
			"service_id": s.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, s)
		c.Next()
	}
}
