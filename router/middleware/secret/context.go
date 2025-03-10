// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "secret"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the secret associated with this context.
func FromContext(c context.Context) *api.Secret {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	r, ok := value.(*api.Secret)
	if !ok {
		return nil
	}

	return r
}

// ToContext adds the hook to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *api.Secret) {
	c.Set(key, s)
}
