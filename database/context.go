// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

// FromContext returns the database Service associated with this context.
func FromContext(c context.Context) Service {
	v := c.Value(key)
	if v == nil {
		return nil
	}

	d, ok := v.(Service)
	if !ok {
		return nil
	}

	return d
}

// ToContext adds the database Service to this context if it supports
// the Setter interface.
func ToContext(c Setter, d Service) {
	c.Set(key, d)
}
