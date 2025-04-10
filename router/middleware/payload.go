// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

// Payload is a middleware function that captures the user provided json body
// and attaches it to the context of every http.Request to be logged.
func Payload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// bind JSON payload from request to be added to context
		var payload any
		_ = c.BindJSON(&payload)

		body, _ := json.Marshal(&payload)

		c.Set("payload", payload)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()
	}
}
