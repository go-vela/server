// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListHooksForRepo gets a list of hooks by repo ID from the database.
func (e *Engine) ListHooksForRepo(ctx context.Context, r *api.Repo, page, perPage int) ([]*api.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing hooks for repo %s", r.GetFullName())

	// variables to store query results and return value
	h := new([]types.Hook)
	hooks := []*api.Hook{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableHook).
		Preload("Build").
		Where("repo_id = ?", r.GetID()).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
		Find(&h).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, hook := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := hook

		result := tmp.ToAPI()
		result.SetRepo(r)

		hooks = append(hooks, result)
	}

	return hooks, nil
}
