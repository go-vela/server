// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"testing"

	"github.com/go-vela/server/queue/models"
)

func TestRedis_Push(t *testing.T) {
	// setup types
	// use global variables in redis_test.go
	_item := &models.Item{
		Build: _build,
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
		redis   *Client
		id      int64
	}{
		{
			failure: false,
			redis:   _redis,
			id:      _item.Build.GetID(),
		},
		{
			failure: true,
			redis:   badItem,
			id:      0,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.redis.Push(context.Background(), "vela", test.id)

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
