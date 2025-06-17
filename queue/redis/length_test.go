// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/server/queue/models"
)

func TestRedis_Length(t *testing.T) {
	// setup types
	// use global variables in redis_test.go
	_item := &models.Item{
		Build: _build,
	}

	// setup queue item
	bytes, err := json.Marshal(_item)
	if err != nil {
		t.Errorf("unable to marshal queue item: %v", err)
	}

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
			want:   4,
		},
		{
			routes: []string{"vela", "vela:second", "phony"},
			want:   6,
		},
	}

	// run tests
	for _, test := range tests {
		for _, route := range test.routes {
			err := _redis.Push(context.Background(), route, bytes)
			if err != nil {
				t.Errorf("unable to push item to queue: %v", err)
			}
		}
		got, err := _redis.Length(context.Background())
		if err != nil {
			t.Errorf("Length returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Length is %v, want %v", got, test.want)
		}
	}
}
