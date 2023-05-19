// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
