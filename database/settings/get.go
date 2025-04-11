// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/database/types"
)

// GetSettings gets platform settings from the database.
func (e *Engine) GetSettings(ctx context.Context) (*settings.Platform, error) {
	e.logger.Trace("getting platform settings")

	// variable to store query results
	s := new(types.Platform)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(TableSettings).
		Where("id = ?", 1).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// return the settings
	return s.ToAPI(), nil
}
