// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// AllowlistSchedule is a middleware function that attaches the allowlistschedule used
// to limit which repos can utilize the schedule feature within the system.
func AllowlistSchedule(allowlistschedule []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("allowlistschedule", allowlistschedule)
		c.Next()
	}
}
