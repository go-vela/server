// SPDX-License-Identifier: Apache-2.0

package native

import (
	api "github.com/go-vela/server/api/types"
)

// SetSettings sets the api settings type.
func (c *client) GetSettings() *api.CompilerSettings {
	return c.CompilerSettings
}

// SetSettings sets the api settings type.
func (c *client) SetSettings(s *api.Settings) {
	if s != nil {
		c.SetCloneImage(s.GetCloneImage())
		c.SetTemplateDepth(s.GetTemplateDepth())
		c.SetStarlarkExecLimit(s.GetStarlarkExecLimit())
	}
}
