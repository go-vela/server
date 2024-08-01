// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "pipeline"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Pipeline associated with this context.
func FromContext(c context.Context) *library.Pipeline {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	b, ok := value.(*library.Pipeline)
	if !ok {
		return nil
	}

	return b
}

// ToContext adds the Pipeline to this context if it supports
// the Setter interface.
func ToContext(c Setter, b *library.Pipeline) {
	c.Set(key, b)
}
