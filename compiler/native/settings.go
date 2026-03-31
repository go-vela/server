// SPDX-License-Identifier: Apache-2.0

package native

import (
	"github.com/go-vela/server/api/types/settings"
)

// GetSettings retrieves the api settings type.
func (c *Client) GetSettings() settings.Compiler {
	return c.Compiler
}

// SetSettings sets the api settings type.
func (c *Client) SetSettings(s *settings.Platform) {
	if s != nil {
		c.SetCloneImage(s.GetCloneImage())
		c.SetTemplateDepth(s.GetTemplateDepth())
		c.SetStarlarkExecLimit(s.GetStarlarkExecLimit())

		// copy pointer fields directly to preserve nil vs empty-slice distinction
		if s.Compiler != nil {
			c.Compiler.BlockedImages = s.Compiler.BlockedImages
			c.Compiler.WarnImages = s.Compiler.WarnImages
		}
	}
}
