// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code in update.go
package dashboard

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/sirupsen/logrus"
)

// CreateUser creates a new user in the database.
func (e *engine) CreateDashboard(ctx context.Context, d *api.Dashboard) (*api.Dashboard, error) {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetName(),
	}).Tracef("creating dashboard %s in the database", d.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#UserFromLibrary
	dashboard := FromAPI(d)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Validate
	err := dashboard.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableDashboard).Create(dashboard)

	return dashboard.ToAPI(), result.Error
}
