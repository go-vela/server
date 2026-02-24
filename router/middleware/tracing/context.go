// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/go-vela/server/tracing"
)

const key = "tracing"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(any, any)
}

// FromContext returns the associated value with this context.
func FromContext(c context.Context) *tracing.Client {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	tc, ok := value.(*tracing.Client)
	if !ok {
		return nil
	}

	return tc
}

// ToContext adds the value to this context if it supports
// the Setter interface.
func ToContext(c Setter, tc *tracing.Client) {
	c.Set(key, tc)
}
