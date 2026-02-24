// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "pipeline"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Pipeline associated with this context.
func FromContext(c context.Context) *api.Pipeline {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	b, ok := value.(*api.Pipeline)
	if !ok {
		return nil
	}

	return b
}

// ToContext adds the Pipeline to this context if it supports
// the Setter interface.
func ToContext(c Setter, b *api.Pipeline) {
	c.Set(key, b)
}
