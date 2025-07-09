// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"github.com/gin-gonic/gin"
)

// AppWebhookSecret is a middleware function that attaches the Vela GH app secret used for
// validating incoming app install webhooks to the context of every http.Request.
func AppWebhookSecret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app-webhook-secret", secret)
		c.Next()
	}
}
