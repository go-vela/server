// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestCache_New(t *testing.T) {
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
				Driver:          "redis",
				Address:         fmt.Sprintf("redis://%s", _redis.Addr()),
				InstallTokenKey: "example",
				Cluster:         false,
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:          "unsupported cache",
				Address:         "cache://cache.example.com",
				InstallTokenKey: "example",
				Cluster:         false,
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
