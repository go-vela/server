// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// DeleteDashboard deletes an existing dashboard from the database.
func (e *engine) DeleteDashboard(ctx context.Context, d *api.Dashboard) error {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Tracef("deleting dashboard %s from the database", d.GetID())

	dashboard := FromAPI(d)

	// send query to the database
	return e.client.
		Table(constants.TableDashboard).
		Delete(dashboard).
		Error
}
