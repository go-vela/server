// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/go-vela/server/constants"
)

// EvictInstallTokens evicts the installation tokens from Redis.
func (c *Client) EvictInstallToken(ctx context.Context, token string) error {
	// compute the HMAC used as the Redis key suffix
	h := hmac.New(sha256.New, []byte(c.config.InstallTokenKey))

	h.Write([]byte(token))

	hmacHex := hex.EncodeToString(h.Sum(nil))

	key := constants.CacheInstallTokenPrefix + hmacHex

	// delete the key
	return c.Redis.Unlink(ctx, key).Err()
}

// EvictBuildInstallTokens evicts all installation tokens associated with a build from Redis.
func (c *Client) EvictBuildInstallTokens(ctx context.Context, build int64) error {
	indexKey := fmt.Sprintf("%s%d", constants.CacheBuildIndexPrefix, build)

	keys, err := c.Redis.SMembers(ctx, indexKey).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		if err := c.Redis.Unlink(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	if err := c.Redis.Del(ctx, indexKey).Err(); err != nil {
		return err
	}

	return nil
}

// EvictInstallStatusToken evicts the installation status token from Redis.
func (c *Client) EvictInstallStatusToken(ctx context.Context, build int64) error {
	key := fmt.Sprintf("%s%d", constants.CacheInstallStatusTokenPrefix, build)

	// delete the key
	return c.Redis.Unlink(ctx, key).Err()
}
