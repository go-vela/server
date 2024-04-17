// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"github.com/go-vela/server/api/types/settings"
)

// SetSettings sets the api settings type in the Engine.
func (c *client) GetSettings() settings.Queue {
	return c.Queue
}

// SetSettings sets the api settings type in the Engine.
func (c *client) SetSettings(s *settings.Platform) {
	if s != nil {
		c.SetRoutes(s.GetRoutes())
	}
}
