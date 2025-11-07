// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	yml "go.yaml.in/yaml/v3"

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
		// TODO:
		// we should be able to use c.YAML here from gin,
		// but there's some incompatibility with us creating yaml.Node
		// with the go.yaml.in/yaml/v3 package and gin using gopkg.in/yaml.v3
		// when calling c.YAML. When gin switches to go.yaml.in/yaml/v3 we can
		// switch to using c.YAML here.
		body, err := yml.Marshal(value)
		if err != nil {
			reason := fmt.Errorf("unable to marshal YAML response: %w", err)
			util.HandleError(c, http.StatusInternalServerError, reason)

			return
		}

		c.Data(http.StatusOK, gin.MIMEYAML, body)
	}
}
