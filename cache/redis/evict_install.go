// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (c *Client) EvictInstallToken(ctx context.Context, token string) error {
	// compute the HMAC used as the Redis key suffix
	h := hmac.New(sha256.New, []byte(c.config.InstallTokenKey))

	h.Write([]byte(token))

	hmacHex := hex.EncodeToString(h.Sum(nil))

	key := "install_token:" + hmacHex

	// delete the key
	return c.Redis.Del(ctx, key).Err()
}
