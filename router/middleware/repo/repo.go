// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/util"
)

// Retrieve gets the repo in the given context.
func Retrieve(c *gin.Context) *api.Repo {
	return FromContext(c)
}

// Establish sets the repo in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		o := org.Retrieve(c)
		ctx := c.Request.Context()

		rParam := util.PathParameter(c, "repo")
		if len(rParam) == 0 {
			retErr := fmt.Errorf("no repo parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		l.Debugf("reading repo %s", rParam)

		// construct full name
		fullName := fmt.Sprintf("%s/%s", o, rParam)

		r, err := database.FromContext(c).GetRepoForOrg(ctx, fullName)
		if err != nil {
			retErr := fmt.Errorf("unable to read repo %s: %w", fullName, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"repo":    r.GetName(),
			"repo_id": r.GetID(),
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, r)
		c.Next()
	}
}
