// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"

	api "github.com/go-vela/server/api/types"
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

	StoreCheckRuns(ctx context.Context, buildID int64, checkRuns []models.CheckRun, timeout int32) error
	StoreStepCheckRuns(ctx context.Context, stepID int64, checkRuns []models.CheckRun, timeout int32) error
	GetCheckRuns(ctx context.Context, build *api.Build) ([]models.CheckRun, error)
	GetStepCheckRuns(ctx context.Context, step *api.Step) ([]models.CheckRun, error)
}
