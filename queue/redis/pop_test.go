// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-vela/types"
	"gopkg.in/square/go-jose.v2/json"
)

func TestRedis_Pop(t *testing.T) {
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
	_redis, err := NewTest("vela")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// push item to queue
	err = _redis.Redis.RPush(context.Background(), "vela", bytes).Err()
	if err != nil {
		t.Errorf("unable to push item to queue: %v", err)
	}

	// setup timeout redis mock
	timeout, err := NewTest("vela")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}
	// overwrite timeout to be 1s
	timeout.config.Timeout = 1 * time.Second

	// setup badChannel redis mock
	badChannel, err := NewTest("vela")
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}
	// overwrite channel to be invalid
	badChannel.config.Channels = nil

	// push something to badChannel queue
	err = badChannel.Redis.RPush(context.Background(), "vela", bytes).Err()
	if err != nil {
		t.Errorf("unable to push item to queue: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		redis   *client
		want    *types.Item
	}{
		{
			failure: false,
			redis:   _redis,
			want:    _item,
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
		got, err := test.redis.Pop(context.Background(), nil)

		if test.failure {
			if err == nil {
				t.Errorf("Pop should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Pop returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Pop is %v, want %v", got, test.want)
		}
	}
}
