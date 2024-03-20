// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetHookForBuild gets a hook by build ID from the database.
func (e *engine) GetHookForBuild(_ context.Context, b *library.Build) (*library.Hook, error) {
	e.logger.Tracef("getting hook for build %d from the database", b.GetID())

	// variable to store query results
	h := new(database.Hook)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Where("build_id = ?", b.GetID()).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	// return the hook
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Hook.ToLibrary
	return h.ToLibrary(), nil
}
