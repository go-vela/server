// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestRedis_New(t *testing.T) {
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
		address string
	}{
		{
			failure: false,
			address: fmt.Sprintf("redis://%s", _redis.Addr()),
		},
		{
			failure: true,
			address: "",
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			context.Background(),
			WithAddress(test.address),
			WithInstallTokenKey("installKey"),
			WithCluster(false),
		)

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
