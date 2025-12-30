// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"

	"github.com/go-vela/server/cache/models"
)

// Service represents the interface for Vela integrating
// with the different supported Queue backends.
type Service interface {
	// Service Interface Functions

	// Driver defines a function that outputs
	// the configured queue driver.
	Driver() string

	StoreInstallToken(ctx context.Context, token *models.InstallToken, timeout int32) error
	GetInstallToken(ctx context.Context, token string) (*models.InstallToken, error)
	EvictInstallToken(ctx context.Context, token string) error
}
