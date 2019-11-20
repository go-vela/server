// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
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

	buf := bytes.Buffer{}

	if w.Empty() {
		return constants.DefaultRoute, nil
	}

	// Build route
	if !strings.EqualFold(strings.ToLower(w.Flavor), "") {
		buf.WriteString(fmt.Sprintf(":%s", w.Flavor))
	}
	if !strings.EqualFold(strings.ToLower(w.Platform), "") {
		buf.WriteString(fmt.Sprintf(":%s", w.Platform))
	}

	return strings.TrimLeft(buf.String(), ":"), nil
}
