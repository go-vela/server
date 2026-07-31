// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"

	"github.com/go-vela/server/cache/models"
)

// StoreTemplateContents stores it in Redis with a TTL.
func (c *Client) StoreTemplateContents(ctx context.Context, key string, t *models.TemplateEntry) error {
	metaBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	// store a small marker value (or metadata JSON if needed)
	err = c.Redis.Set(ctx, key, metaBytes, c.config.TemplateTTL).Err()
	if err != nil {
		return err
	}

	return nil
}
