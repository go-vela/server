// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateDashboard updates an existing dashboard in the database.
func (e *engine) UpdateDashboard(ctx context.Context, d *api.Dashboard) (*api.Dashboard, error) {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Tracef("creating dashboard %s", d.GetID())

	dashboard := types.DashboardFromAPI(d)

	err := dashboard.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.Table(constants.TableDashboard).Save(dashboard).Error
	if err != nil {
		return nil, err
	}

	return dashboard.ToAPI(), nil
}
