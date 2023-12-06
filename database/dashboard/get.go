// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetDashboard gets a dashboard by ID from the database.
func (e *engine) GetDashboard(ctx context.Context, id string) (*library.Dashboard, error) {
	e.logger.Tracef("getting dashboard %s from the database", id)

	// variable to store query results
	r := new(database.Dashboard)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableDashboard).
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, err
	}

	// return the decrypted dashboard
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Dashboard.ToLibrary
	return r.ToLibrary(), nil
}
