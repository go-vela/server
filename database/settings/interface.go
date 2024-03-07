// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
)

// SettingsInterface represents the Vela interface for settings
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type SettingsInterface interface {
	// CreateSettings defines a function that creates a platform settings record.
	CreateSettings(context.Context, *string) (*string, error)
	// GetSettings defines a function that gets platform settings.
	GetSettings(context.Context) (*string, error)
	// UpdateSettings defines a function that updates platform settings.
	UpdateSettings(context.Context, *string) (*string, error)
}
