// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "dashboard"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, any)
}

// FromContext returns the Dashboard associated with this context.
func FromContext(c context.Context) *api.Dashboard {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	b, ok := value.(*api.Dashboard)
	if !ok {
		return nil
	}

	return b
}

// ToContext adds the Dashboard to this context if it supports
// the Setter interface.
func ToContext(c Setter, b *api.Dashboard) {
	c.Set(key, b)
}
