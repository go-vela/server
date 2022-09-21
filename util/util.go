// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package util

import (
	"html"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
)

// HandleError appends the error to the handler chain for logging and outputs it.
func HandleError(c *gin.Context, status int, err error) {
	msg := err.Error()
	//nolint:errcheck // ignore checking error
	c.Error(err)
	c.AbortWithStatusJSON(status, types.Error{Message: &msg})
}

// MaxInt is a helper function to clamp the integer which
// prevents it from being higher then the provided value.
//
// Currently, Go only supports float64 via math. ( max | min ).
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt is a helper function to clamp the integer which
// prevents it from being lower then the provided value.
//
// Currently, Go only supports float64 via math. ( max | min ).
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// FormParameter safely captures a form parameter from the context
// by removing any new lines and HTML escaping the value.
func FormParameter(c *gin.Context, parameter string) string {
	return EscapeValue(c.Request.FormValue(parameter))
}

// QueryParameter safely captures a query parameter from the context
// by removing any new lines and HTML escaping the value.
func QueryParameter(c *gin.Context, parameter, value string) string {
	return EscapeValue(c.DefaultQuery(parameter, value))
}

// PathParameter safely captures a path parameter from the context
// by removing any new lines and HTML escaping the value.
func PathParameter(c *gin.Context, parameter string) string {
	return EscapeValue(c.Param(parameter))
}

// EscapeValue safely escapes any string by removing any new lines and HTML escaping it.
func EscapeValue(value string) string {
	// replace all new lines in the value
	escaped := strings.Replace(strings.Replace(value, "\n", "", -1), "\r", "", -1)

	// HTML escape the new line escaped value
	return html.EscapeString(escaped)
}
