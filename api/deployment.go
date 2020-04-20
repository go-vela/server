// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateDeployment represents the API handler to
// create a deployment in the configured backend.
func CreateDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// GetDeployments represents the API handler to capture
// a list of deployments from the configured backend.
func GetDeployments(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// GetDeployment represents the API handler to
// capture a deployment from the configured backend.
func GetDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// UpdateDeployment represents the API handler to
// update a deployment in the configured backend.
func UpdateDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}

// DeleteDeployment represents the API handler to
// remove a deployment from the configured backend.
func DeleteDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}
