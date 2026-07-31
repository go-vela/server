// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
)

// ExtendTemplateExpiry extends the TTL for a template entry in Redis.
func (c *Client) ExtendTemplateExpiry(ctx context.Context, key string) error {
	// extend the TTL for the template entry
	err := c.Redis.Expire(ctx, key, c.config.TemplateTTL).Err()
	if err != nil {
		return err
	}

	return nil
}
