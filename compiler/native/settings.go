// SPDX-License-Identifier: Apache-2.0

package native

import (
	api "github.com/go-vela/server/api/types"
)

// UpdateFromSettings sets the api settings type in the Engine.
func (c *client) UpdateFromSettings(s *api.Settings) {
	if s != nil {
		c.CloneImage = s.GetCloneImage()
		c.TemplateDepth = int(s.GetTemplateDepth())
		c.StarlarkExecLimit = s.GetStarlarkExecLimit()
	}
}
