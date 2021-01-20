// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"context"
)

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the secret Service
// associated with this context.
func FromContext(c context.Context, key string) Service {
	// get secret value from context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast secret value to expected Service type
	s, ok := v.(Service)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the secret Service to this
// context if it supports the Setter interface.
func ToContext(c Setter, key string, s Service) {
	c.Set(key, s)
}
