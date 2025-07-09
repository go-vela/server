// SPDX-License-Identifier: Apache-2.0

package github

import (
	"github.com/go-vela/server/api/types/settings"
)

// GetSettings retrieves the api settings type in the Engine.
func (c *Client) GetSettings() settings.SCM {
	return c.SCM
}

// SetSettings sets the api settings type in the Engine.
func (c *Client) SetSettings(s *settings.Platform) {
	if s != nil {
		c.SetRepoRoleMap(s.GetRepoRoleMap())
		c.SetOrgRoleMap(s.GetOrgRoleMap())
		c.SetTeamRoleMap(s.GetTeamRoleMap())
	}
}
