// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "schedule"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(any, any)
}

// FromContext returns the Schedule associated with this context.
func FromContext(c context.Context) *api.Schedule {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*api.Schedule)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Schedule to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *api.Schedule) {
	c.Set(key, s)
}
