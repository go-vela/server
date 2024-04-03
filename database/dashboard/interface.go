// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// DashboardInterface represents the Vela interface for repo
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type DashboardInterface interface {
	// Dashboard Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateDashboardTable defines a function that creates the dashboards table.
	CreateDashboardTable(context.Context, string) error

	// Dashboard Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateDashboard defines a function that creates a dashboard.
	CreateDashboard(context.Context, *api.Dashboard) (*api.Dashboard, error)
	// DeleteDashboard defines a function that deletes a dashboard.
	DeleteDashboard(context.Context, *api.Dashboard) error
	// GetDashboard defines a function that gets a dashboard by ID.
	GetDashboard(context.Context, string) (*api.Dashboard, error)
	// UpdateDashboard defines a function that updates a dashboard.
	UpdateDashboard(context.Context, *api.Dashboard) (*api.Dashboard, error)
}
