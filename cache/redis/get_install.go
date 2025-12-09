// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/go-vela/server/cache/models"
)

func (c *Client) GetInstallToken(ctx context.Context, token string) (*models.InstallToken, error) {
	h := hmac.New(sha256.New, []byte(c.config.InstallTokenKey))

	h.Write([]byte(token))

	hmacHex := hex.EncodeToString(h.Sum(nil))

	key := "install_token:" + hmacHex

	meta, err := c.Redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	installToken := new(models.InstallToken)

	err = json.Unmarshal(meta, installToken)
	if err != nil {
		return nil, err
	}

	installToken.Token = token

	return installToken, nil
}
