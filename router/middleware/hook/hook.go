// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// Retrieve gets the hook in the given context.
func Retrieve(c *gin.Context) *library.Hook {
	return FromContext(c)
}

// Establish sets the hook in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		ctx := c.Request.Context()

		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", o, util.PathParameter(c, "repo"))
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		hParam := util.PathParameter(c, "hook")
		if len(hParam) == 0 {
			retErr := fmt.Errorf("no hook parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		number, err := strconv.Atoi(hParam)
		if err != nil {
			retErr := fmt.Errorf("malformed hook parameter provided: %s", hParam)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		l.Debugf("reading hook %s/%d", r.GetFullName(), number)

		h, err := database.FromContext(c).GetHookForRepo(ctx, r, number)
		if err != nil {
			retErr := fmt.Errorf("unable to read hook %s/%d: %w", r.GetFullName(), number, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"hook": h.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, h)
		c.Next()
	}
}
