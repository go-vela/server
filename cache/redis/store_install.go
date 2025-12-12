// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
)

// StoreInstallToken computes an HMAC-SHA256 of the token and stores it in Redis with a TTL.
func (c *Client) StoreInstallToken(ctx context.Context, t *models.InstallToken, repo *api.Repo) error {
	meta := new(models.InstallToken)
	meta.Repositories = t.Repositories
	meta.Permissions = t.Permissions
	meta.Expiration = t.Expiration

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	// set TTL based on repo timeout
	ttl := time.Minute * time.Duration(repo.GetTimeout())

	h := hmac.New(sha256.New, []byte(c.config.InstallTokenKey))

	h.Write([]byte(t.Token))

	hmacHex := hex.EncodeToString(h.Sum(nil))

	key := "install_token:" + hmacHex

	// store a small marker value (or metadata JSON if needed)
	err = c.Redis.Set(ctx, key, metaBytes, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
