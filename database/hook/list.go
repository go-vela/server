// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListHooks gets a list of all hooks from the database.
func (e *Engine) ListHooks(ctx context.Context) ([]*api.Hook, error) {
	e.logger.Trace("listing all hooks")

	// variables to store query results and return value
	h := new([]types.Hook)
	hooks := []*api.Hook{}

	// send query to the database and store result in variable
	err := e.client.
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
		err = hook.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo for hook %d: %v", hook.ID.Int64, err)
		}

		hooks = append(hooks, hook.ToAPI())
	}

	return hooks, nil
}
