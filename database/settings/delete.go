// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
)

// DeleteSettings removes platform settings from the database.
func (e *engine) DeleteSettings(ctx context.Context, s *settings.Platform) error {
	e.logger.Trace("deleting platform settings from the database")

	worker := FromAPI(s)

	// send query to the database
	return e.client.
		Table(TableSettings).
		Delete(worker).
		Error
}
