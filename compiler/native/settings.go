// SPDX-License-Identifier: Apache-2.0

package native

import (
	"github.com/go-vela/server/api/types/settings"
)

// GetSettings retrieves the api settings type.
func (c *client) GetSettings() settings.Compiler {
	return c.Compiler
}

// SetSettings sets the api settings type.
func (c *client) SetSettings(s *settings.Platform) {
	if s != nil {
		c.SetCloneImage(s.GetCloneImage())
		c.SetTemplateDepth(s.GetTemplateDepth())
		c.SetStarlarkExecLimit(s.GetStarlarkExecLimit())
	}
}
