// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AllDeployments represents the API handler to
// captures all deployments stored in the database.
func AllDeployments(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// UpdateDeployment represents the API handler to
// update any deployment stored in the database.
func UpdateDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}
