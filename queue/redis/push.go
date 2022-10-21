// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"
	"errors"

	"golang.org/x/crypto/nacl/sign"
)

// Push inserts an item to the specified channel in the queue.
func (c *client) Push(ctx context.Context, channel string, item []byte) error {
	c.Logger.Tracef("pushing item to queue %s", channel)

	var signed []byte
	var out []byte

	// check for this on startup
	if c.config.SigningPrivateKey == nil || len(*c.config.SigningPrivateKey) != 64 {
		return errors.New("no valid signing private key provided")
	}

	c.Logger.Tracef("signing item for queue %s", channel)

	// sign the item using the private key generated using sign
	//
	// https://pkg.go.dev/golang.org/x/crypto@v0.1.0/nacl/sign
	signed = sign.Sign(out, item, c.config.SigningPrivateKey)

	// build a redis queue command to push an item to queue
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#Client.RPush
	pushCmd := c.Redis.RPush(ctx, channel, signed)

	// blocking call to push an item to queue and return err
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#IntCmd.Err
	err := pushCmd.Err()
	if err != nil {
		return err
	}

	return nil
}
