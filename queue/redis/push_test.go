// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-vela/types"
)

func TestRedis_Push(t *testing.T) {
	// setup types
	// use global variables in redis_test.go
	_item := &types.Item{
		Build: _build,
		Repo:  _repo,
		User:  _user,
	}

	// setup queue item
	_bytes, err := json.Marshal(_item)
	if err != nil {
		t.Errorf("unable to marshal queue item: %v", err)
	}

	// setup redis mock
	_redis, err := NewTest(_signingPrivateKey, _signingPublicKey, "vela")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup redis mock
	badItem, err := NewTest(_signingPrivateKey, _signingPublicKey, "vela")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		redis   *client
		bytes   []byte
	}{
		{
			failure: false,
			redis:   _redis,
			bytes:   _bytes,
		},
		{
			failure: true,
			redis:   badItem,
			bytes:   nil,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.redis.Push(context.Background(), "vela", test.bytes)

		if test.failure {
			if err == nil {
				t.Errorf("Push should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Push returned err: %v", err)
		}
	}
}
