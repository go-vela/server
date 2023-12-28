// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountHooks gets the count of all hooks from the database.
func (e *engine) CountHooks(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all hooks from the database")

	// variable to store query results
	var h int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Count(&h).
		Error

	return h, err
}
