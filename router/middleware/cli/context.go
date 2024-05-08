// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"context"

	"github.com/urfave/cli/v2"
)

const key = "cli"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the cli context associated with this context.
func FromContext(c context.Context) *cli.Context {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*cli.Context)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the cli context to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *cli.Context) {
	c.Set(key, s)
}
