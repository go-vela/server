// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"context"
)

// key defines the key type for storing
// the queue Service in the context.
const key = "queue"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the queue Service
// associated with this context.
func FromContext(c context.Context) Service {
	// get queue value from context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast queue value to expected Service type
	q, ok := v.(Service)
	if !ok {
		return nil
	}

	return q
}

// ToContext adds the queue Service to this
// context if it supports the Setter interface.
func ToContext(c Setter, q Service) {
	c.Set(key, q)
}
