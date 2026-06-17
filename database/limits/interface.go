// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// LimitInterface represents the Vela interface for limit
// functions with the supported Database backends.
type LimitInterface interface {
	// Limit Data Definition Language Functions

	// CreateOrgBuildLimitTable defines a function that creates the org_build_limits table.
	CreateOrgBuildLimitTable(context.Context, string) error

	// Limit Data Manipulation Language Functions

	// CreateOrgBuildLimit defines a function that creates an org build limit.
	CreateOrgBuildLimit(context.Context, *api.OrgBuildLimit) (*api.OrgBuildLimit, error)
	// GetOrgBuildLimit defines a function that gets an org build limit by org.
	GetOrgBuildLimit(context.Context, string) (*api.OrgBuildLimit, error)
	// UpdateOrgBuildLimit defines a function that updates an org build limit.
	UpdateOrgBuildLimit(context.Context, *api.OrgBuildLimit) (*api.OrgBuildLimit, error)
}
