// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetBuildForRepo gets a build by repo ID and number from the database.
func (e *engine) GetBuildForRepo(ctx context.Context, r *api.Repo, number int) (*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"build": number,
		"org":   r.GetOrg(),
		"repo":  r.GetName(),
	}).Tracef("getting build %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	b := new(types.Build)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("repo_id = ?", r.GetID()).
		Where("number = ?", number).
		Take(b).
		Error
	if err != nil {
		return nil, err
	}

	err = b.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo %s/%s: %v", r.GetOrg(), r.GetName(), err)

		return b.ToAPI(), nil
	}

	return b.ToAPI(), nil
}
