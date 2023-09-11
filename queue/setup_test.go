// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestQueue_Setup_Redis(t *testing.T) {
	// setup types
	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}
	defer _redis.Close()

	_setup := &Setup{
		Driver:  "redis",
		Address: fmt.Sprintf("redis://%s", _redis.Addr()),
		Routes:  []string{"foo"},
		Cluster: false,
	}

	_, err = _setup.Redis()
	if err != nil {
		t.Errorf("Redis returned err: %v", err)
	}
}

func TestQueue_Setup_Kafka(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:  "kafka",
		Address: "kafka://kafka.example.com",
		Routes:  []string{"foo"},
		Cluster: false,
	}

	got, err := _setup.Kafka()
	if err == nil {
		t.Errorf("Kafka should have returned err")
	}

	if got != nil {
		t.Errorf("Kafka is %v, want nil", got)
	}
}

func TestQueue_Setup_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Driver:     "redis",
				Address:    "redis://redis.example.com",
				Routes:     []string{"foo"},
				Cluster:    false,
				PrivateKey: "sdafhdsaohf",
				PublicKey:  "adwfwo",
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:     "kafka",
				Address:    "kafka://kafka.example.com",
				Routes:     []string{"foo"},
				Cluster:    false,
				PrivateKey: "sdafhdsaohf",
				PublicKey:  "adwfwo",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:  "redis",
				Address: "redis://redis.example.com/",
				Routes:  []string{"foo"},
				Cluster: false,
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:  "redis",
				Address: "redis.example.com",
				Routes:  []string{"foo"},
				Cluster: false,
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:  "",
				Address: "redis://redis.example.com",
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
		{
			failure: true,
			setup: &Setup{
				Driver:  "redis",
				Address: "redis://redis.example.com",
				Routes:  []string{},
				Cluster: false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.setup.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
