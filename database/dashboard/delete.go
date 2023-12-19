// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteDashboard deletes an existing dashboard from the database.
func (e *engine) DeleteDashboard(ctx context.Context, d *library.Dashboard) error {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Tracef("deleting dashboard %s from the database", d.GetID())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#DashboardFromLibrary
	dashboard := database.DashboardFromLibrary(d)

	// send query to the database
	return e.client.
		Table(constants.TableDashboard).
		Delete(dashboard).
		Error
}
