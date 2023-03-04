// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// DefaultRepoEvents is a middleware function that attaches the defaultRepoEvents
// to enable the server to override the default repo event.
func DefaultRepoEvents(defaultRepoEvents []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("defaultRepoEvents", defaultRepoEvents)
		c.Next()
	}
}
