// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "build"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Build associated with this context.
func FromContext(c context.Context) *library.Build {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	b, ok := value.(*library.Build)
	if !ok {
		return nil
	}

	return b
}

// ToContext adds the Build to this context if it supports
// the Setter interface.
func ToContext(c Setter, b *library.Build) {
	c.Set(key, b)
}
