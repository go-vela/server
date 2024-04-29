// SPDX-License-Identifier: Apache-2.0

package keyset

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetKeySet gets a key by UUID from the database if active.
func (e *engine) GetActiveKeySet(ctx context.Context, id string) (api.JWK, error) {
	e.logger.Tracef("getting key %s from the database", id)

	// variable to store query results
	j := new(types.JWK)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableKeySet).
		Where("id = ?", id).
		Where("active = ?", true).
		Take(j).
		Error
	if err != nil {
		return j.ToAPI(), err
	}

	return j.ToAPI(), nil
}
