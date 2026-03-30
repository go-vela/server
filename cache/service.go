// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
)

// Service represents the interface for Vela integrating
// with the different supported Queue backends.
type Service interface {
	// Service Interface Functions

	// Driver defines a function that outputs
	// the configured queue driver.
	Driver() string

	StoreInstallToken(ctx context.Context, token string, build int64, timeout int32) error
	GetInstallToken(ctx context.Context, token string) error
	EvictInstallToken(ctx context.Context, token string) error
	EvictBuildInstallTokens(ctx context.Context, build int64) error

	StoreInstallStatusToken(ctx context.Context, build int64, token string) error
	GetInstallStatusToken(ctx context.Context, build int64) (string, error)
	EvictInstallStatusToken(ctx context.Context, build int64) error

	StorePermissionToken(ctx context.Context, installID int64, token string) error
	GetPermissionToken(ctx context.Context, installID int64) (string, error)
}
