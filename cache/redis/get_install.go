// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

func (c *Client) GetInstallToken(ctx context.Context, token string) error {
	h := hmac.New(sha256.New, []byte(c.config.InstallTokenKey))

	h.Write([]byte(token))

	hmacHex := hex.EncodeToString(h.Sum(nil))

	key := constants.CacheInstallTokenPrefix + hmacHex

	_, err := c.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

// GetInstallStatusToken retrieves the installation status token from Redis.
func (c *Client) GetInstallStatusToken(ctx context.Context, build int64) (string, error) {
	key := fmt.Sprintf("%s%d", constants.CacheInstallStatusTokenPrefix, build)

	token, err := c.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	decoded, err := util.Decrypt(c.config.InstallTokenKey, []byte(token))
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

// GetPermissionToken retrieves the permission token from Redis.
func (c *Client) GetPermissionToken(ctx context.Context, installID int64) (string, error) {
	key := fmt.Sprintf("%s%d", constants.CachePermissionTokenPrefix, installID)

	token, err := c.Redis.Get(ctx, key).Result()
	if err != nil {
		//nolint:nilerr // cache miss return non-error with empty token
		return "", nil
	}

	decoded, err := util.Decrypt(c.config.InstallTokenKey, []byte(token))
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
