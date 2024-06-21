// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// Retrieve gets the schedule in the given context.
func Retrieve(c *gin.Context) *api.Schedule {
	return FromContext(c)
}

// Establish sets the schedule in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		r := repo.Retrieve(c)
		ctx := c.Request.Context()

		sParam := util.PathParameter(c, "schedule")
		if len(sParam) == 0 {
			retErr := fmt.Errorf("no schedule parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		l.Debugf("reading schedule %s", sParam)

		s, err := database.FromContext(c).GetScheduleForRepo(ctx, r, sParam)
		if err != nil {
			retErr := fmt.Errorf("unable to read schedule %s for repo %s: %w", sParam, r.GetFullName(), err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"schedule":    s.GetName(),
			"schedule_id": s.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, s)
		c.Next()
	}
}
