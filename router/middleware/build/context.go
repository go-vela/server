// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "build"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Build associated with this context.
func FromContext(c context.Context) *api.Build {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	b, ok := value.(*api.Build)
	if !ok {
		return nil
	}

	return b
}

// ToContext adds the Build to this context if it supports
// the Setter interface.
func ToContext(c Setter, b *api.Build) {
	c.Set(key, b)
}
