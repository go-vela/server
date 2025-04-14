// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/database/types"
)

// UpdateSettings updates a platform settings in the database.
func (e *Engine) UpdateSettings(ctx context.Context, s *settings.Platform) (*settings.Platform, error) {
	e.logger.Trace("updating platform settings in the database")

	// cast the api type to database type
	dbS := types.SettingsFromAPI(s)

	// validate the necessary fields are populated
	err := dbS.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(TableSettings).
		Save(dbS.Nullify()).Error
	if err != nil {
		return nil, err
	}

	s = dbS.ToAPI()

	return s, nil
}
