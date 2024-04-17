// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CreateDashboard creates a new dashboard in the database.
func (e *engine) CreateDashboard(ctx context.Context, d *api.Dashboard) (*api.Dashboard, error) {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetName(),
	}).Tracef("creating dashboard %s in the database", d.GetName())

	dashboard := FromAPI(d)

	err := dashboard.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableDashboard).Create(dashboard)

	return dashboard.ToAPI(), result.Error
}
