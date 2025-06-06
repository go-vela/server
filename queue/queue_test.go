// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestQueue_New(t *testing.T) {
	// setup types
	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}
	defer _redis.Close()

	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Driver:     "redis",
				Address:    fmt.Sprintf("redis://%s", _redis.Addr()),
				Routes:     []string{"foo"},
				Cluster:    false,
				PrivateKey: "bOiFT7Y9e0jpOqaapTa3NzUkAve3VdRvyowgsY/vtlcK5L4RADOh9uTe1UVLdu3l/a0hvhiIkkLidUwVBhASWA==",
				PublicKey:  "CuS+EQAzofbk3tVFS3bt5f2tIb4YiJJC4nVMFQYQElg=",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:  "kafka",
				Address: "kafka://kafka.example.com",
				Routes:  []string{"foo"},
				Cluster: false,
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:  "pubsub",
				Address: "pubsub://pubsub.example.com",
				Routes:  []string{"foo"},
				Cluster: false,
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:  "redis",
				Address: "",
				Routes:  []string{"foo"},
				Cluster: false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(context.Background(), test.setup)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}
