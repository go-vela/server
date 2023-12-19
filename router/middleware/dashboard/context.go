// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "dashboard"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Dashboard associated with this context.
func FromContext(c context.Context) *library.Dashboard {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	b, ok := value.(*library.Dashboard)
	if !ok {
		return nil
	}

	return b
}

// ToContext adds the Dashboard to this context if it supports
// the Setter interface.
func ToContext(c Setter, b *library.Dashboard) {
	c.Set(key, b)
}
