// SPDX-License-Identifier: Apache-2.0

package redis

import (
	api "github.com/go-vela/server/api/types"
)

// SetSettings sets the api settings type in the Engine.
func (c *client) GetSettings() *api.QueueSettings {
	return c.QueueSettings
}

// SetSettings sets the api settings type in the Engine.
func (c *client) SetSettings(s *api.Settings) {
	if s != nil {
		c.SetQueueRoutes(s.GetQueueRoutes())
	}
}
