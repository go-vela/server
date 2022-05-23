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
	// nolint: errcheck // ignore checking error
	c.Error(err)
	c.AbortWithStatusJSON(status, types.Error{Message: &msg})
}

// GetParameter safely captures a parameter from the context by
// removing any new lines and HTML escaping the value.
func GetParameter(c *gin.Context, parameter string) string {
	// capture the raw value for the parameter
	raw := c.Param(parameter)

	// replace all new lines in the value for the parameter
	escaped := strings.Replace(strings.Replace(raw, "\n", "", -1), "\r", "", -1)

	// HTML escape the new line escaped value for the parameter
	return html.EscapeString(escaped)
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
