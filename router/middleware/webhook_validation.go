// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"github.com/gin-gonic/gin"
)

// WebhookValidation determines whether or not incoming webhooks are validated coming from Github
// This is primarily intended for local development.
func WebhookValidation(validate bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("webhookvalidation", validate)
		c.Next()
	}
}
