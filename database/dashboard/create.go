// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateDashboard creates a new dashboard in the database.
func (e *Engine) CreateDashboard(ctx context.Context, d *api.Dashboard) (*api.Dashboard, error) {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetName(),
	}).Tracef("creating dashboard %s", d.GetName())

	dashboard := types.DashboardFromAPI(d)

	err := dashboard.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableDashboard).
		Create(dashboard)

	return dashboard.ToAPI(), result.Error
}
