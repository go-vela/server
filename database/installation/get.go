// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetInstallation gets an installation by ID from the database.
func (e *Engine) GetInstallation(ctx context.Context, target string) (*api.Installation, error) {
	e.logger.Tracef("getting installation %s", target)

	// variable to store query results
	i := new(types.Installation)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableInstallation).
		Where("target = ?", target).
		Take(i).
		Error
	if err != nil {
		return nil, err
	}

	// return the installation
	return i.ToAPI(), nil
}
