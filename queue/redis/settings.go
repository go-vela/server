// SPDX-License-Identifier: Apache-2.0

package redis

import (
	api "github.com/go-vela/server/api/types"
)

// UpdateFromSettings takes settings and updates the queue.
func (c *client) UpdateFromSettings(s *api.Settings) {
	c.Logger.Trace("updating queue using settings")

	c.config.Channels = s.GetQueueRoutes()
}
