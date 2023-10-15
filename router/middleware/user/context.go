// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "user"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the User associated with this context.
func FromContext(c context.Context) *library.User {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	u, ok := value.(*library.User)
	if !ok {
		return nil
	}

	return u
}

// ToContext adds the User to this context if it supports
// the Setter interface.
func ToContext(c Setter, u *library.User) {
	c.Set(key, u)
}
