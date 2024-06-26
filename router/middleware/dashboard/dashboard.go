// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// Retrieve gets the build in the given context.
func Retrieve(c *gin.Context) *api.Dashboard {
	return FromContext(c)
}

// Establish sets the build in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		u := user.Retrieve(c)
		ctx := c.Request.Context()

		id := util.PathParameter(c, "dashboard")
		if len(id) == 0 {
			userBoards := u.GetDashboards()
			if len(userBoards) == 0 {
				retErr := fmt.Errorf("no available dashboards for user %s", u.GetName())
				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			id = userBoards[0]
		}

		l.Debugf("reading dashboard %s", id)

		d, err := database.FromContext(c).GetDashboard(ctx, id)
		if err != nil {
			retErr := fmt.Errorf("unable to read dashboard %s: %w", id, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"dashboard": d.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, d)
		c.Next()
	}
}
