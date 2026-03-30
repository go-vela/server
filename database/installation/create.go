// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateInstallation creates a new installation in the database.
func (e *Engine) CreateInstallation(ctx context.Context, i *api.Installation) (*api.Installation, error) {
	e.logger.WithFields(logrus.Fields{
		"installation": i.GetInstallID(),
	}).Tracef("creating installation %d", i.GetInstallID())

	// cast the API type to database type
	installation := types.InstallationFromAPI(i)

	// validate the necessary fields are populated
	err := installation.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableInstallation).
		Create(installation)

	return installation.ToAPI(), result.Error
}
