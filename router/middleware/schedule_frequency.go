// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

// ScheduleFrequency is a middleware function that attaches the scheduleminimumfrequency used
// to limit the frequency which schedules can be run within the system.
func ScheduleFrequency(scheduleFrequency time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("scheduleminimumfrequency", scheduleFrequency)
		c.Next()
	}
}
