// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the worker in the given context.
func Retrieve(c *gin.Context) *api.Worker {
	return FromContext(c)
}

// Establish sets the worker in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		wParam := util.PathParameter(c, "worker")
		if len(wParam) == 0 {
			retErr := fmt.Errorf("no worker parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		logrus.Debugf("Reading worker %s", wParam)

		w, err := database.FromContext(c).GetWorkerForHostname(ctx, wParam)
		if err != nil {
			retErr := fmt.Errorf("unable to read worker %s: %w", wParam, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		ToContext(c, w)
		c.Next()
	}
}
