// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package org

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/util"
)

// Retrieve gets the org in the given context.
func Retrieve(c *gin.Context) string {
	return FromContext(c)
}

// Establish used to check if org param is used only.
func Establish() gin.HandlerFunc {
	return func(c *gin.Context) {
		oParam := html.EscapeString(c.Param("org"))
		if len(oParam) == 0 {
			retErr := fmt.Errorf("no org parameter provided")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		ToContext(c, oParam)

		c.Next()
	}
}
