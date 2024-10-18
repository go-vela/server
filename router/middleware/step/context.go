// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "step"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Step associated with this context.
func FromContext(c context.Context) *api.Step {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*api.Step)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Step to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *api.Step) {
	c.Set(key, s)
}
