// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the build in the given context.
func Retrieve(c *gin.Context) *library.Dashboard {
	return FromContext(c)
}

// Establish sets the build in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"dashboard": id,
			"user":      u.GetName(),
		}).Debugf("reading dashboard %s", id)

		d, err := database.FromContext(c).GetDashboard(ctx, id)
		if err != nil {
			retErr := fmt.Errorf("unable to read dashboard %s: %w", id, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		ToContext(c, d)
		c.Next()
	}
}
