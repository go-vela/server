// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package util

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/types"
)

// HandleError appends the error to the handler chain for logging and outputs it
func HandleError(c *gin.Context, status int, err error) {
	msg := err.Error()
	c.Error(err)
	c.AbortWithStatusJSON(status, types.Error{Message: &msg})
}

// Helper functions to clamp integers - go only supports float64 via math.(max|min)
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
