// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/go-vela/types"
	"golang.org/x/crypto/nacl/sign"
)

// Pop grabs an item from the specified channel off the queue.
func (c *client) Pop(ctx context.Context) (*types.Item, error) {
	c.Logger.Tracef("popping item from queue %s", c.config.Channels)

	// build a redis queue command to pop an item from queue
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#Client.BLPop
	popCmd := c.Redis.BLPop(ctx, c.config.Timeout, c.config.Channels...)

	// blocking call to pop item from queue
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#StringSliceCmd.Result
	result, err := popCmd.Result()
	if err != nil {
		// BLPOP timeout
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	// check for this on startup
	if c.config.SigningPublicKey == nil || len(*c.config.SigningPublicKey) != 32 {
		return nil, errors.New("no valid signing public key provided")
	}

	_item := []byte(result[1])
	var opened []byte
	var out []byte

	// open item
	opened, ok := sign.Open(out, _item, c.config.SigningPublicKey)
	if !ok {
		return nil, errors.New("unable to open queue item")
	}

	_item = opened

	item := new(types.Item)

	// unmarshal result into queue item
	err = json.Unmarshal(_item, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
