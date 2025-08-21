// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "testattachment"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the TestAttachment associated with this context.
func FromContext(c context.Context) *api.TestAttachment {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	ta, ok := value.(*api.TestAttachment)
	if !ok {
		return nil
	}

	return ta
}

// ToContext adds the TestAttachment to this context if it supports
// the Setter interface.
func ToContext(c Setter, ta *api.TestAttachment) {
	c.Set(key, ta)
}
