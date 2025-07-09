// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"github.com/go-vela/server/api/types/settings"
)

// GetSettings retrieves the api settings type in the Engine.
func (c *Client) GetSettings() settings.Queue {
	return c.Queue
}

// SetSettings sets the api settings type in the Engine.
func (c *Client) SetSettings(s *settings.Platform) {
	if s != nil {
		c.SetRoutes(s.GetRoutes())
	}
}
