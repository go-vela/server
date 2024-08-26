// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetDashboard gets a dashboard by UUID from the database.
func (e *engine) GetDashboard(ctx context.Context, id string) (*api.Dashboard, error) {
	e.logger.Tracef("getting dashboard %s", id)

	// variable to store query results
	d := new(types.Dashboard)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableDashboard).
		Where("id = ?", id).
		Take(d).
		Error
	if err != nil {
		return nil, err
	}

	return d.ToAPI(), nil
}
