// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
)

// SettingsInterface represents the Vela interface for settings
// functions with the supported Database backends.
//

type SettingsInterface interface {
	// CreateSettings defines a function that creates a platform settings record.
	CreateSettings(context.Context, *settings.Platform) (*settings.Platform, error)
	// GetSettings defines a function that gets platform settings.
	GetSettings(context.Context) (*settings.Platform, error)
	// UpdateSettings defines a function that updates platform settings.
	UpdateSettings(context.Context, *settings.Platform) (*settings.Platform, error)
}
