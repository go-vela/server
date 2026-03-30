// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteInstallation deletes an existing installation from the database.
func (e *Engine) DeleteInstallation(ctx context.Context, i *api.Installation) error {
	e.logger.WithFields(logrus.Fields{
		"installation": i.GetInstallID(),
	}).Tracef("deleting installation %d", i.GetInstallID())

	// cast the API type to database type
	installation := types.InstallationFromAPI(i)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableInstallation).
		Where("install_id = ?", installation.InstallID.Int64).
		Delete(&types.Installation{}).
		Error
}
