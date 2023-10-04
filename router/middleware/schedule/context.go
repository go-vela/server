// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "schedule"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Schedule associated with this context.
func FromContext(c context.Context) *library.Schedule {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*library.Schedule)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Schedule to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *library.Schedule) {
	c.Set(key, s)
}
