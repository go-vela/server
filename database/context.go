// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"context"
)

const key = "database"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the database Interface associated with this context.
func FromContext(c context.Context) Interface {
	v := c.Value(key)
	if v == nil {
		return nil
	}

	d, ok := v.(Interface)
	if !ok {
		return nil
	}

	return d
}

// ToContext adds the database Interface to this context if it supports
// the Setter interface.
func ToContext(c Setter, d Interface) {
	c.Set(key, d)
}
