// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "artifact"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Artifact associated with this context.
func FromContext(c context.Context) *api.Artifact {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	a, ok := value.(*api.Artifact)
	if !ok {
		return nil
	}

	return a
}

// ToContext adds the Artifact to this context if it supports
// the Setter interface.
func ToContext(c Setter, a *api.Artifact) {
	c.Set(key, a)
}
