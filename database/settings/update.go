// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
)

// UpdateSettings updates a platform settings in the database.
func (e *engine) UpdateSettings(_ context.Context, s *settings.Platform) (*settings.Platform, error) {
	e.logger.Trace("updating platform settings in the database")

	// cast the api type to database type
	settings := FromAPI(s)

	// validate the necessary fields are populated
	err := settings.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.Table(TableSettings).Save(settings.Nullify()).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}
