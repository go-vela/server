// SPDX-License-Identifier: Apache-2.0

package step

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// Retrieve gets the step in the given context.
func Retrieve(c *gin.Context) *api.Step {
	return FromContext(c)
}

// Establish sets the step in the given context.
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

		sParam := util.PathParameter(c, "step")
		if len(sParam) == 0 {
			retErr := fmt.Errorf("no step parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		number, err := strconv.Atoi(sParam)
		if err != nil {
			retErr := fmt.Errorf("malformed step parameter provided: %s", sParam)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		l.Debugf("reading step %d", number)

		s, err := database.FromContext(c).GetStepForBuild(ctx, b, number)
		if err != nil {
			retErr := fmt.Errorf("unable to read step %s/%d/%d: %w", r.GetFullName(), b.GetNumber(), number, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"step":    s.GetNumber(),
			"step_id": s.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, s)
		c.Next()
	}
}
