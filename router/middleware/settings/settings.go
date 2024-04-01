// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// Retrieve gets the settings in the given context.
func Retrieve(c *gin.Context) *api.Settings {
	return FromContext(c)
}

// Establish sets the settings in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		logrus.Debug("Reading settings")

		s, err := database.FromContext(c).GetSettings(ctx)
		if err != nil {
			retErr := fmt.Errorf("unable to read settings: %w", err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		ToContext(c, s)
		c.Next()
	}
}
