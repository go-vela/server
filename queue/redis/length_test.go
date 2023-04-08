// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-vela/types"
	"gopkg.in/square/go-jose.v2/json"
)

func TestRedis_Length(t *testing.T) {
	// setup types
	// use global variables in redis_test.go
	_item := &types.Item{
		Build:    _build,
		Pipeline: _steps,
		Repo:     _repo,
		User:     _user,
	}

	// setup queue item
	bytes, err := json.Marshal(_item)
	if err != nil {
		t.Errorf("unable to marshal queue item: %v", err)
	}

	// setup redis mock
	_redis, err := NewTest("vela", "vela:second", "vela:third")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// setup tests
	tests := []struct {
		channels []string
		want     int64
	}{
		{
			channels: []string{"vela"},
			want:     1,
		},
		{
			channels: []string{"vela", "vela:second", "vela:third"},
			want:     4,
		},
		{
			channels: []string{"vela", "vela:second", "phony"},
			want:     6,
		},
	}

	// run tests
	for _, test := range tests {
		for _, channel := range test.channels {
			err := _redis.Push(context.Background(), channel, bytes)
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
