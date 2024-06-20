// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetHookForRepo gets a hook by repo ID and number from the database.
func (e *engine) GetHookForRepo(ctx context.Context, r *api.Repo, number int) (*api.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"hook": number,
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting hook %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	h := new(types.Hook)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Preload("Repo").
		Preload("Repo.Owner").
		Preload("Build").
		Where("repo_id = ?", r.GetID()).
		Where("number = ?", number).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	err = h.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo %s/%s: %v", r.GetOrg(), r.GetName(), err)
	}

	return h.ToAPI(), nil
}
