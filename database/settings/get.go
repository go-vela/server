// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// GetSettings gets platform settings from the database.
func (e *engine) GetSettings(ctx context.Context) (*api.Settings, error) {
	e.logger.Trace("getting platform settings from the database")

	// variable to store query results
	s := new(Settings)

	// send query to the database and store result in variable
	err := e.client.
		Table(TableSettings).
		// todo: how to ensure this is always a singleton at the first row
		Where("id = ?", 1).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// return the settings
	return s.ToAPI(), nil
}
