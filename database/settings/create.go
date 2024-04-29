// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
)

// CreateSettings creates a platform settings record in the database.
func (e *engine) CreateSettings(ctx context.Context, s *settings.Platform) (*settings.Platform, error) {
	e.logger.Tracef("creating platform settings in the database with %v", s.String())

	// cast the api type to database type
	settings := FromAPI(s)

	// validate the necessary fields are populated
	err := settings.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.Table(TableSettings).Create(settings).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}
