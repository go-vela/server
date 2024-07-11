// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "hook"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the hook associated with this context.
func FromContext(c context.Context) *api.Hook {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	r, ok := value.(*api.Hook)
	if !ok {
		return nil
	}

	return r
}

// ToContext adds the hook to this context if it supports
// the Setter interface.
func ToContext(c Setter, h *api.Hook) {
	c.Set(key, h)
}
