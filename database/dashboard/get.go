// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// GetDashboard gets a dashboard by UUID from the database.
func (e *engine) GetDashboard(ctx context.Context, id string) (*api.Dashboard, error) {
	e.logger.Tracef("getting dashboard %s from the database", id)

	// variable to store query results
	r := new(Dashboard)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableDashboard).
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, err
	}

	return r.ToAPI(), nil
}
