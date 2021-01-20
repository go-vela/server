// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "worker"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Worker associated with this context.
func FromContext(c context.Context) *library.Worker {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	w, ok := value.(*library.Worker)
	if !ok {
		return nil
	}

	return w
}

// ToContext adds the Worker to this context if it supports
// the Setter interface.
func ToContext(c Setter, w *library.Worker) {
	c.Set(key, w)
}
