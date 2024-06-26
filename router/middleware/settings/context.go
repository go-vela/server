// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
)

const key = "settings"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Settings associated with this context.
func FromContext(c context.Context) *settings.Platform {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*settings.Platform)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Settings to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *settings.Platform) {
	c.Set(key, s)
}
