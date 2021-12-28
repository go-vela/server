// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
)

// Route decides which route a build gets placed within the queue.
func (c *client) Route(w *pipeline.Worker) (string, error) {
	c.Logger.Tracef("deciding route from queue channels %s", c.config.Channels)

	// create buffer to store route
	buf := bytes.Buffer{}

	// if pipline does not specify route information return default
	//
	// https://github.com/go-vela/types/blob/master/constants/queue.go#L10
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

	return strings.TrimLeft(buf.String(), ":"), nil
}
