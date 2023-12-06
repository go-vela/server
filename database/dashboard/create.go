// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code in update.go
package dashboard

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateUser creates a new user in the database.
func (e *engine) CreateDashboard(ctx context.Context, d *library.Dashboard) (*library.Dashboard, error) {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetName(),
	}).Tracef("creating dashboard %s in the database", d.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#UserFromLibrary
	dashboard := database.DashboardFromLibrary(d)

	dashboard.ID = uuid.New()

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Validate
	err := dashboard.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableDashboard).Create(dashboard)

	return dashboard.ToLibrary(), result.Error
}
