// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListHooks gets a list of all hooks from the database.
func (e *engine) ListHooks(ctx context.Context) ([]*api.Hook, error) {
	e.logger.Trace("listing all hooks")

	// variables to store query results and return value
	count := int64(0)
	h := new([]types.Hook)
	hooks := []*api.Hook{}

	// count the results
	count, err := e.CountHooks(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return hooks, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		WithContext(ctx).
		Table(constants.TableHook).
		Preload("Repo").
		Preload("Repo.Owner").
		Preload("Build").
		Find(&h).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, hook := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := hook

		err = tmp.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo for hook %d: %v", tmp.ID.Int64, err)
		}

		hooks = append(hooks, tmp.ToAPI())
	}

	return hooks, nil
}
