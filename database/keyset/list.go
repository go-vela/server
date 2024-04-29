// SPDX-License-Identifier: Apache-2.0

package keyset

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListKeySets gets a list of all schedules from the database.
func (e *engine) ListKeySets(ctx context.Context) ([]api.JWK, error) {
	e.logger.Trace("listing all keysets from the database")

	k := new([]types.JWK)
	keys := []api.JWK{}

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableKeySet).
		Find(&k).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, key := range *k {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := key

		// convert query result to API type
		keys = append(keys, tmp.ToAPI())
	}

	return keys, nil
}
