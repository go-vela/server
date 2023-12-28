// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// ScheduleFrequency is a middleware function that attaches the scheduleminimumfrequency used
// to limit the frequency which schedules can be run within the system.
func ScheduleFrequency(scheduleFrequency time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("scheduleminimumfrequency", scheduleFrequency)
		c.Next()
	}
}
