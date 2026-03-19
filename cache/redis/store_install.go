// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-vela/server/cache/models"
	"github.com/go-vela/server/constants"
)

// StoreInstallToken computes an HMAC-SHA256 of the token and stores it in Redis with a TTL.
func (c *Client) StoreInstallToken(ctx context.Context, t *models.InstallToken, build int64, timeout int32) error {
	meta := new(models.InstallToken)
	meta.InstallID = t.InstallID
	meta.Repositories = t.Repositories
	meta.Permissions = t.Permissions
	meta.Expiration = t.Expiration
	meta.Timeout = timeout

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	// set TTL based on repo timeout
	ttl := time.Minute * time.Duration(timeout)

	h := hmac.New(sha256.New, []byte(c.config.InstallTokenKey))

	h.Write([]byte(t.Token))

	hmacHex := hex.EncodeToString(h.Sum(nil))

	key := constants.CacheInstallTokenPrefix + hmacHex

	// store a small marker value (or metadata JSON if needed)
	err = c.Redis.Set(ctx, key, metaBytes, ttl).Err()
	if err != nil {
		return err
	}

	// add the key to a Redis set for the build to enable eviction of all tokens for a build
	idxKey := fmt.Sprintf("%s%d", constants.CacheBuildIndexPrefix, build)

	err = c.Redis.SAdd(ctx, idxKey, key).Err()
	if err != nil {
		return err
	}

	// set an expiry on the index in case eviction routine is never called
	err = c.Redis.Expire(ctx, idxKey, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

// StoreInstallStatusToken stores the installation status token in Redis with a TTL.
func (c *Client) StoreInstallStatusToken(ctx context.Context, build int64, token string) error {
	// set TTL to 59 minutes to ensure it expires before the GitHub token does
	ttl := time.Minute * 59

	key := fmt.Sprintf("%s%d", constants.CacheInstallStatusTokenPrefix, build)

	// store the token with a TTL
	err := c.Redis.Set(ctx, key, token, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
