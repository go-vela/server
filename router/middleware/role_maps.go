// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// RoleMaps is a middleware function that attaches the role maps
// to the context of every http.Request.
func RoleMaps(rMap, oMap, tMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("repo.roles-map", rMap)
		c.Set("org.roles-map", oMap)
		c.Set("team.roles-map", tMap)

		c.Next()
	}
}
