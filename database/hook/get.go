// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetHook gets a hook by ID from the database.
func (e *engine) GetHook(ctx context.Context, id int64) (*api.Hook, error) {
	e.logger.Tracef("getting hook %d from the database", id)

	// variable to store query results
	h := new(types.Hook)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Preload("Repo").
		Preload("Build").
		Where("id = ?", id).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	return h.ToAPI(), nil
}
