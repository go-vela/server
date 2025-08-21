// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// Retrieve gets the test attachment in the given context.
func Retrieve(c *gin.Context) *api.TestAttachment {
	return FromContext(c)
}

// Establish sets the test attachment in the given context.
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

		l.Debugf("reading test attachment")

		ta, err := database.FromContext(c).GetTestAttachmentForBuild(ctx, b)
		if err != nil {
			retErr := fmt.Errorf("unable to read test attachment %s/%d/reports/testreport/attachment/: %w", r.GetFullName(), b.GetNumber(), err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"testreport_id": ta.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, ta)
		c.Next()
	}
}
