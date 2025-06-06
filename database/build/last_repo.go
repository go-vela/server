// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// LastBuildForRepo gets the last build by repo ID and branch from the database.
func (e *Engine) LastBuildForRepo(ctx context.Context, r *api.Repo, branch string) (*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting last build for repo %s", r.GetFullName())

	// variable to store query results
	b := new(types.Build)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where("branch = ?", branch).
		Order("number DESC").
		Take(b).
		Error
	if err != nil {
		// check if the query returned a record not found error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// the record will not exist if it is a new repo
			return nil, nil
		}

		return nil, err
	}

	result := b.ToAPI()
	result.SetRepo(r)

	return result, nil
}
