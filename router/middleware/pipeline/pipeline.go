// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the pipeline in the given context.
func Retrieve(c *gin.Context) *library.Pipeline {
	return FromContext(c)
}

// Establish sets the pipeline in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)
		ctx := c.Request.Context()

		if r == nil {
			retErr := fmt.Errorf("repo %s/%s not found", util.PathParameter(c, "org"), util.PathParameter(c, "repo"))

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		p := util.PathParameter(c, "pipeline")
		if len(p) == 0 {
			retErr := fmt.Errorf("no pipeline parameter provided")

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		entry := fmt.Sprintf("%s/%s", r.GetFullName(), p)

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"org":      o,
			"pipeline": p,
			"repo":     r.GetName(),
			"user":     u.GetName(),
		}).Debugf("reading pipeline %s", entry)

		pipeline, err := database.FromContext(c).GetPipelineForRepo(ctx, p, r)
		if err != nil { // assume the pipeline doesn't exist in the database yet (before pipeline support was added)
			// send API call to capture the pipeline configuration file
			config, err := scm.FromContext(c).ConfigBackoff(ctx, u, r, p)
			if err != nil {
				retErr := fmt.Errorf("unable to get pipeline configuration for %s: %w", entry, err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}

			// parse and compile the pipeline configuration file
			_, pipeline, err = compiler.FromContext(c).
				Duplicate().
				WithCommit(p).
				WithMetadata(c.MustGet("metadata").(*types.Metadata)).
				WithRepo(r).
				WithUser(u).
				Compile(config)
			if err != nil {
				retErr := fmt.Errorf("unable to compile pipeline configuration for %s: %w", entry, err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}

		ToContext(c, pipeline)

		c.Next()
	}
}
