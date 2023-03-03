// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "initstep"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the InitStep associated with this context.
func FromContext(c context.Context) *library.InitStep {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*library.InitStep)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the InitStep to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *library.InitStep) {
	c.Set(key, s)
}
