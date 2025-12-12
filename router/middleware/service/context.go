// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "service"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(any, any)
}

// FromContext returns the Service associated with this context.
func FromContext(c context.Context) *api.Service {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*api.Service)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Service to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *api.Service) {
	c.Set(key, s)
}
