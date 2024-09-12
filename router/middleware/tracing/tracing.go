// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/tracing"
)

// Retrieve gets the value in the given context.
func Retrieve(c *gin.Context) *tracing.Client {
	return FromContext(c)
}

// Establish sets the value in the given context.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		tc := Retrieve(c)

		l.Debugf("reading tracing client from context")

		ToContext(c, tc)
		c.Next()
	}
}
