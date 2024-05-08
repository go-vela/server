// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	cliMiddleware "github.com/go-vela/server/router/middleware/cli"
)

// CLI is a middleware function that attaches the urface cli client
// to the context of every http.Request.
func CLI(cliCtx *cli.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		cliMiddleware.ToContext(c, cliCtx)

		c.Next()
	}
}
