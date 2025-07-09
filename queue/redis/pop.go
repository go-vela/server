// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/nacl/sign"

	"github.com/go-vela/server/queue/models"
)

// Pop grabs an item from the specified route off the queue.
func (c *Client) Pop(ctx context.Context, inRoutes []string) (*models.Item, error) {
	// define routes to pop from
	var routes []string

	// if routes were supplied, use those
	if len(inRoutes) > 0 {
		routes = inRoutes
	} else {
		routes = c.GetRoutes()
	}

	c.Logger.Tracef("popping item from queue %s", routes)

	// build a redis queue command to pop an item from queue
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#Client.BLPop
	popCmd := c.Redis.BLPop(ctx, c.config.Timeout, routes...)

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

	// extract signed item from pop results
	signed := []byte(result[1])

	var opened, out []byte

	// open the item using the public key generated using sign
	//
	// https://pkg.go.dev/golang.org/x/crypto@v0.1.0/nacl/sign
	opened, ok := sign.Open(out, signed, c.config.PublicKey)
	if !ok {
		return nil, errors.New("unable to open signed item")
	}

	// unmarshal result into queue item
	item := new(models.Item)

	err = json.Unmarshal(opened, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
