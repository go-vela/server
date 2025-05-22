// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"testing"
)

func TestRedis_RouteLength(t *testing.T) {
	// setup redis mock
	_redis, err := NewTest(_signingPrivateKey, _signingPublicKey, "vela", "vela:second", "vela:third")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		routes []string
		want   int64
	}{
		{
			routes: []string{"vela"},
			want:   1,
		},
		{
			routes: []string{"vela", "vela:second", "vela:third"},
			want:   2,
		},
		{
			routes: []string{"vela", "vela:second", "phony"},
			want:   3,
		},
	}

	// run tests
	id := int64(0)
	for _, test := range tests {
		for _, channel := range test.routes {
			id++
			err := _redis.Push(context.Background(), channel, id)
			if err != nil {
				t.Errorf("unable to push item to queue: %v", err)
			}
		}
		got, err := _redis.RouteLength(context.Background(), "vela")
		if err != nil {
			t.Errorf("RouteLength returned err: %v", err)
		}

		if got != test.want {
			t.Errorf("Length is %v, want %v", got, test.want)
		}
	}
}
