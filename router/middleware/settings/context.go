// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "settings"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Settings associated with this context.
func FromContext(c context.Context) *api.Settings {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*api.Settings)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Settings to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *api.Settings) {
	c.Set(key, s)
}
