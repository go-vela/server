// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetHookByWebhookID gets a single hook with a matching webhook id in the database.
func (e *engine) GetHookByWebhookID(ctx context.Context, webhookID int64) (*api.Hook, error) {
	e.logger.Tracef("getting a hook with webhook id %d from the database", webhookID)

	// variable to store query results
	h := new(types.Hook)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Preload("Repo").
		Preload("Repo.Owner").
		Preload("Build").
		Where("webhook_id = ?", webhookID).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	err = h.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo for hook %d: %v", h.ID.Int64, err)
	}

	return h.ToAPI(), nil
}
