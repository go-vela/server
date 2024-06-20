// SPDX-License-Identifier: Apache-2.0

package org

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/util"
)

// Retrieve gets the org in the given context.
func Retrieve(c *gin.Context) string {
	return FromContext(c)
}

// Establish used to check if org param is used only.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)

		oParam := util.PathParameter(c, "org")

		if len(oParam) == 0 {
			retErr := fmt.Errorf("no org parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		l = l.WithFields(logrus.Fields{
			"org": oParam,
		})

		// update the logger with the new fields
		c.Set("logger", l)

		ToContext(c, oParam)
		c.Next()
	}
}
