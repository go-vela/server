// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// Payload is a middleware function that captures the user provided json body
// and attaches it to the context of every http.Request to be logged
func Payload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// bind JSON payload from request to be added to context
		var payload interface{}
		_ = c.BindJSON(&payload)

		body, _ := json.Marshal(&payload)

		c.Set("payload", payload)

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		c.Next()
	}
}
