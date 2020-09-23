// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package org

import (
	"context"
)

const key = "org"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Repo associated with this context.
func FromContext(c context.Context) string {
	value := c.Value(key)
	if value == nil {
		return ""
	}

	o, ok := value.(string)
	if !ok {
		return ""
	}

	return o
}

// ToContext adds the Repo to this context if it supports
// the Setter interface.
func ToContext(c Setter, o string) {
	c.Set(key, o)
}
