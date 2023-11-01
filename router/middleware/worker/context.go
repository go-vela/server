// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "worker"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Worker associated with this context.
func FromContext(c context.Context) *api.Worker {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	w, ok := value.(*api.Worker)
	if !ok {
		return nil
	}

	return w
}

// ToContext adds the Worker to this context if it supports
// the Setter interface.
func ToContext(c Setter, w *api.Worker) {
	c.Set(key, w)
}
