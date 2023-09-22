// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"errors"

	"golang.org/x/crypto/nacl/sign"
)

// Push inserts an item to the specified channel in the queue.
func (c *client) Push(ctx context.Context, channel string, item []byte) error {
	c.Logger.Tracef("pushing item to queue %s", channel)

	// ensure the item to be pushed is valid
	// go-redis RPush does not support nil as of v9.0.2
	//
	// https://github.com/redis/go-redis/pull/1960
	if item == nil {
		return errors.New("item is nil")
	}

	var signed []byte

	var out []byte

	// this should already be validated on startup
	if c.config.PrivateKey == nil || len(*c.config.PrivateKey) != 64 {
		return errors.New("no valid signing private key provided")
	}

	c.Logger.Tracef("signing item for queue %s", channel)

	// sign the item using the private key generated using sign
	//
	// https://pkg.go.dev/golang.org/x/crypto@v0.1.0/nacl/sign
	signed = sign.Sign(out, item, c.config.PrivateKey)

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
