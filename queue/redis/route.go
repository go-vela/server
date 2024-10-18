// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// Route decides which route a build gets placed within the queue.
func (c *client) Route(w *pipeline.Worker) (string, error) {
	c.Logger.Tracef("deciding route from queue routes %s", c.GetRoutes())

	// create buffer to store route
	buf := bytes.Buffer{}

	// if pipline does not specify route information return default
	//
	// https://github.com/go-vela/types/blob/main/constants/queue.go#L10
	if w.Empty() {
		return constants.DefaultRoute, nil
	}

	// append flavor to route
	if !strings.EqualFold(strings.ToLower(w.Flavor), "") {
		buf.WriteString(fmt.Sprintf(":%s", w.Flavor))
	}

	// append platform to route
	if !strings.EqualFold(strings.ToLower(w.Platform), "") {
		buf.WriteString(fmt.Sprintf(":%s", w.Platform))
	}

	route := strings.TrimLeft(buf.String(), ":")

	for _, r := range c.GetRoutes() {
		if strings.EqualFold(route, r) {
			return route, nil
		}
	}

	return "", fmt.Errorf("invalid route %s provided", route)
}
