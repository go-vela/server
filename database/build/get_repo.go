// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetBuildForRepo gets a build by repo ID and number from the database.
func (e *engine) GetBuildForRepo(ctx context.Context, r *api.Repo, number int) (*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"build": number,
		"org":   r.GetOrg(),
		"repo":  r.GetName(),
	}).Tracef("getting build %s/%d", r.GetFullName(), number)

	// variable to store query results
	b := new(types.Build)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where("number = ?", number).
		Take(b).
		Error
	if err != nil {
		return nil, err
	}

	result := b.ToAPI()
	result.SetRepo(r)

	return result, nil
}
