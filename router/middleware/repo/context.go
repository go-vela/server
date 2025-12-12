// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "repo"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(any, any)
}

// FromContext returns the Repo associated with this context.
func FromContext(c context.Context) *api.Repo {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	r, ok := value.(*api.Repo)
	if !ok {
		return nil
	}

	return r
}

// ToContext adds the Repo to this context if it supports
// the Setter interface.
func ToContext(c Setter, r *api.Repo) {
	c.Set(key, r)
}
