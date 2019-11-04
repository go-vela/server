// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "step"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Step associated with this context.
func FromContext(c context.Context) *library.Step {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*library.Step)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Step to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *library.Step) {
	c.Set(key, s)
}
