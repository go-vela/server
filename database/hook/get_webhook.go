// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetHookByWebhookID gets a single hook with a matching webhook id in the database.
func (e *engine) GetHookByWebhookID(ctx context.Context, webhookID int64) (*library.Hook, error) {
	e.logger.Tracef("getting a hook with webhook id %d from the database", webhookID)

	// variable to store query results
	h := new(database.Hook)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Where("webhook_id = ?", webhookID).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	// return the hook
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Hook.ToLibrary
	return h.ToLibrary(), nil
}
