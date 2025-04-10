// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"context"

	"github.com/urfave/cli/v3"
)

const key = "cli"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, any)
}

// FromContext returns the cli command associated with this context.
func FromContext(c context.Context) *cli.Command {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*cli.Command)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the cli command to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *cli.Command) {
	c.Set(key, s)
}
