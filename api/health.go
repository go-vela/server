// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health represents the API handler to
// report the health status for Vela.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
