// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetOldestExecutableForRepo gets an executable by repo ID from the database with lowest created_at value.
func (e *engine) GetOldestExecutableForRepo(ctx context.Context, r *library.Repo) (*library.BuildExecutable, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("grabbing oldest executable for repo %s from the database", r.GetFullName())

	// variable to store query results and return value
	executable := new(database.BuildExecutable)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuildExecutable).
		Where("repo_id = ?", r.GetID()).
		Order("created_at ASC").
		Take(&executable).
		Error
	if err != nil {
		return nil, err
	}

	return executable.ToLibrary(), nil
}
