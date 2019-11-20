// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateHook represents the API handler to create
// a webhook in the configured backend.
func CreateHook(c *gin.Context) {
	c.String(http.StatusNotImplemented, "This endpoint is not implemented.")
}

// GetHooks represents the API handler to capture a list
// of webhooks from the configured backend.
func GetHooks(c *gin.Context) {
	c.String(http.StatusNotImplemented, "This endpoint is not implemented.")
}

// GetHook represents the API handler to capture a
// webhook from the configured backend.
func GetHook(c *gin.Context) {
	c.String(http.StatusNotImplemented, "This endpoint is not implemented.")
}

// UpdateHook represents the API handler to update
// a webhook in the configured backend.
func UpdateHook(c *gin.Context) {
	c.String(http.StatusNotImplemented, "This endpoint is not implemented.")
}

// DeleteHook represents the API handler to remove
// a webhook from the configured backend.
func DeleteHook(c *gin.Context) {
	c.String(http.StatusNotImplemented, "This endpoint is not implemented.")
}
