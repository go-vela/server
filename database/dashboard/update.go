// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with create.go
package dashboard

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateDashboard updates an existing dashboard in the database.
func (e *engine) UpdateDashboard(ctx context.Context, d *library.Dashboard) (*library.Dashboard, error) {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Tracef("creating dashboard %s in the database", d.GetID())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#DashboardFromLibrary
	dashboard := database.DashboardFromLibrary(d)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Dashboard.Validate
	err := dashboard.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.Table(constants.TableDashboard).Save(dashboard).Error
	if err != nil {
		return nil, err
	}

	return dashboard.ToLibrary(), nil
}
