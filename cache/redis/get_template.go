// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-vela/server/cache/models"
	"github.com/redis/go-redis/v9"
)

func (c *Client) GetTemplateContents(ctx context.Context, key string) (*models.TemplateEntry, error) {
	meta, err := c.Redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	templateEntry := new(models.TemplateEntry)

	err = json.Unmarshal(meta, templateEntry)
	if err != nil {
		return nil, err
	}

	return templateEntry, nil
}
