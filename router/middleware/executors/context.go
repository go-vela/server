// SPDX-License-Identifier: Apache-2.0

package executors

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "executors"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, any)
}

// FromContext returns the executors associated with this context.
func FromContext(c context.Context) []api.Executor {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	e, ok := value.([]api.Executor)
	if !ok {
		return nil
	}

	return e
}

// ToContext adds the executor to this context if it supports
// the Setter interface.
func ToContext(c Setter, e []api.Executor) {
	c.Set(key, e)
}
