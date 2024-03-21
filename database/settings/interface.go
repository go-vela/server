// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// SettingsInterface represents the Vela interface for settings
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type SettingsInterface interface {
	// CreateSettings defines a function that creates a platform settings record.
	CreateSettings(context.Context, *api.Settings) (*api.Settings, error)
	// GetSettings defines a function that gets platform settings.
	GetSettings(context.Context) (*api.Settings, error)
	// UpdateSettings defines a function that updates platform settings.
	UpdateSettings(context.Context, *api.Settings) (*api.Settings, error)
}
