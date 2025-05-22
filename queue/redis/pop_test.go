// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRedis_Pop(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest(_signingPrivateKey, _signingPublicKey, "vela", "custom")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// push item to queue
	err = _redis.Push(t.Context(), "vela", int64(1))
	if err != nil {
		t.Errorf("unable to push item to queue: %v", err)
	}

	// push item to queue with custom channel
	err = _redis.Push(context.Background(), "custom", int64(2))
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

	// push something to badChannel queue
	err = badChannel.Push(context.Background(), "vela", int64(3))
	if err != nil {
		t.Errorf("unable to push item to queue: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		redis   *Client
		want    int64
		routes  []string
	}{
		{
			failure: false,
			redis:   _redis,
			want:    int64(1),
		},
		{
			failure: false,
			redis:   _redis,
			want:    int64(2),
			routes:  []string{"custom"},
		},
		{
			failure: false,
			redis:   timeout,
			want:    0,
		},
		{
			failure: true,
			redis:   badChannel,
			want:    0,
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
