// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/settings"
)

// Settings is a middleware function that fetches the latest settings and
// attaches to the context of every http.Request.
func Settings(d database.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := d.GetSettings(context.Background())
		if err != nil {
			logrus.WithError(err).Warn("unable to get platform settings")

			return
		}

		settings.ToContext(c, s)

		c.Next()
	}
}
