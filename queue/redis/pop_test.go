// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/nacl/sign"

	"github.com/go-vela/server/queue/models"
)

func TestRedis_Pop(t *testing.T) {
	// setup types
	// use global variables in redis_test.go
	_item := &models.Item{
		Build: _build,
	}

	var signed []byte
	var out []byte

	// setup queue item
	bytes, err := json.Marshal(_item)
	if err != nil {
		t.Errorf("unable to marshal queue item: %v", err)
	}

	// setup redis mock
	_redis, err := NewTest(_signingPrivateKey, _signingPublicKey, "vela", "custom")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	signed = sign.Sign(out, bytes, _redis.config.PrivateKey)

	// push item to queue
	err = _redis.Redis.RPush(context.Background(), "vela", signed).Err()
	if err != nil {
		t.Errorf("unable to push item to queue: %v", err)
	}

	// push item to queue with custom channel
	err = _redis.Redis.RPush(context.Background(), "custom", signed).Err()
	if err != nil {
		t.Errorf("unable to push item to queue: %v", err)
	}

	// setup timeout redis mock
	timeout, err := NewTest(_signingPrivateKey, _signingPublicKey, "vela")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}
	// overwrite timeout to be 1s
	timeout.config.Timeout = 1 * time.Second

	// setup badChannel redis mock
	badChannel, err := NewTest(_signingPrivateKey, _signingPublicKey, "vela")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}
	// overwrite channel to be invalid
	badChannel.SetRoutes(nil)

	signed = sign.Sign(out, bytes, badChannel.config.PrivateKey)

	// push something to badChannel queue
	err = badChannel.Redis.RPush(context.Background(), "vela", signed).Err()
	if err != nil {
		t.Errorf("unable to push item to queue: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		redis   *Client
		want    *models.Item
		routes  []string
	}{
		{
			failure: false,
			redis:   _redis,
			want:    _item,
		},
		{
			failure: false,
			redis:   _redis,
			want:    _item,
			routes:  []string{"custom"},
		},
		{
			failure: false,
			redis:   timeout,
			want:    nil,
		},
		{
			failure: true,
			redis:   badChannel,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.redis.Pop(context.Background(), test.routes)

		if test.failure {
			if err == nil {
				t.Errorf("Pop should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Pop returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("Pop() mismatch (-want +got):\n%s", diff)
		}
	}
}
