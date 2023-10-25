// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"testing"
	"time"
)

func TestRedis_Ping_Good(t *testing.T) {
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}

	defer _redis.Close()

	// setup redis mock
	goodRedis, err := New(
		WithAddress(fmt.Sprintf("redis://%s", _redis.Addr())),
		WithChannels("foo"),
		WithCluster(false),
		WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// run tests
	err = goodRedis.Ping(context.Background())

	if err != nil {
		t.Errorf("Ping returned err: %v", err)
	}
}

func TestRedis_Ping_Bad(t *testing.T) {
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}

	defer _redis.Close()

	// setup redis mock
	badRedis, _ := New(
		WithAddress(fmt.Sprintf("redis://%s", _redis.Addr())),
		WithChannels("foo"),
		WithCluster(false),
		WithTimeout(5*time.Second),
	)
	_redis.SetError("not aiv")
	// run tests
	err = badRedis.Ping(context.Background())
	if err == nil {
		t.Errorf("Ping should have returned err: %v", err)
	}
}
