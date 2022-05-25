// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/util"
)

const (
	outputJSON = "json"
	outputYAML = "yaml"
)

// writeOutput is a helper function to return the provided value to the
// request based off the output query parameter provided. If no output
// query parameter is provided, then YAML is used by default.
func writeOutput(c *gin.Context, value interface{}) {
	output := util.QueryParameter(c, "output", outputYAML)

	// format response body based off output query parameter
	switch strings.ToLower(output) {
	case outputJSON:
		c.JSON(http.StatusOK, value)
	case outputYAML:
		fallthrough
	default:
		c.YAML(http.StatusOK, value)
	}
}
