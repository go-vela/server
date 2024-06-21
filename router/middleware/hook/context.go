// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "hook"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Repo associated with this context.
func FromContext(c context.Context) *library.Hook {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	r, ok := value.(*library.Hook)
	if !ok {
		return nil
	}

	return r
}

// ToContext adds the Repo to this context if it supports
// the Setter interface.
func ToContext(c Setter, r *library.Hook) {
	c.Set(key, r)
}
