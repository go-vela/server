// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateBuildExecutable creates a new build executable in the database.
func (e *Engine) CreateBuildExecutable(ctx context.Context, b *api.BuildExecutable) error {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetBuildID(),
	}).Tracef("creating build executable for build %d in the database", b.GetBuildID())

	// convert API type to database type
	executable := types.BuildExecutableFromAPI(b)

	err := executable.Validate()
	if err != nil {
		return err
	}

	// compress data for the build executable
	err = executable.Compress(e.config.CompressionLevel)
	if err != nil {
		return err
	}

	// encrypt the data field for the build executable
	err = executable.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt build executable for build %d: %w", b.GetBuildID(), err)
	}

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableBuildExecutable).
		Create(executable).
		Error
}
