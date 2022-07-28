// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// FeedbackRepo is a middleware function that attaches the feedbackRepo
// to enable the server to post issues to a feedback repo.
func FeedbackRepo(feedbackRepo string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("feedbackRepo", feedbackRepo)
		c.Next()
	}
}
