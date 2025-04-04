// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// Retrieve gets the build in the given context.
func Retrieve(c *gin.Context) *api.Build {
	return FromContext(c)
}

// Establish sets the build in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		r := repo.Retrieve(c)
		ctx := c.Request.Context()

		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", util.PathParameter(c, "org"), util.PathParameter(c, "repo"))
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		bParam := util.PathParameter(c, "build")
		if len(bParam) == 0 {
			retErr := fmt.Errorf("no build parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		number, err := strconv.ParseInt(bParam, 10, 64)
		if err != nil {
			retErr := fmt.Errorf("invalid build parameter provided: %s", bParam)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		l.Debugf("reading build %d", number)

		b, err := database.FromContext(c).GetBuildForRepo(ctx, r, number)
		if err != nil {
			retErr := fmt.Errorf("unable to read build %s/%d: %w", r.GetFullName(), number, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"build":    b.GetNumber(),
			"build_id": b.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, b)
		c.Next()
	}
}
