// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "user"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(any, any)
}

// FromContext returns the User associated with this context.
func FromContext(c context.Context) *api.User {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	u, ok := value.(*api.User)
	if !ok {
		return nil
	}

	return u
}

// ToContext adds the User to this context if it supports
// the Setter interface.
func ToContext(c Setter, u *api.User) {
	c.Set(key, u)
}
