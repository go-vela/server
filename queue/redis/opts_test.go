// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Bose/minisentinel"
	"github.com/alicebob/miniredis/v2"
)

func TestRedis_ClientOpt_WithAddress(t *testing.T) {
	// setup tests

	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}
	defer _redis.Close()

	tests := []struct {
		failure bool
		address string
		want    string
	}{
		{
			failure: false,
			address: fmt.Sprintf("redis://%s", _redis.Addr()),
			want:    fmt.Sprintf("redis://%s", _redis.Addr()),
		},
		{
			failure: true,
			address: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress(test.address),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithAddress should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Address, test.want) {
			t.Errorf("WithAddress is %v, want %v", _service.config.Address, test.want)
		}
	}
}

func TestRedis_ClientOpt_WithChannels(t *testing.T) {
	// setup tests

	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}
	defer _redis.Close()

	tests := []struct {
		failure  bool
		channels []string
		want     []string
	}{
		{
			failure:  false,
			channels: []string{"foo", "bar"},
			want:     []string{"foo", "bar"},
		},
		{
			failure:  true,
			channels: []string{},
			want:     []string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress(fmt.Sprintf("redis://%s", _redis.Addr())),
			WithChannels(test.channels...),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithChannels should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithChannels returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Channels, test.want) {
			t.Errorf("WithChannels is %v, want %v", _service.config.Channels, test.want)
		}
	}
}

func TestRedis_ClientOpt_WithCluster(t *testing.T) {
	// setup tests

	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_primary, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create primary miniredis instance: %v", err)
	}
	defer _primary.Close()

	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_replica, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create primary miniredis instance: %v", err)
	}
	defer _replica.Close()

	// create a local fake redis cluster
	//
	// https://pkg.go.dev/github.com/Bose/minisentinel#Run
	_cluster, err := minisentinel.Run(_primary, minisentinel.WithReplica(_replica))
	if err != nil {
		t.Errorf("unable to create miniredis cluster: %v", err)
	}
	defer _cluster.Close()

	tests := []struct {
		address string
		cluster bool
		want    bool
	}{
		{
			address: fmt.Sprintf("redis://%s,%s", _cluster.MasterInfo().Name, _cluster.Addr()),
			cluster: true,
			want:    true,
		},
		{
			address: fmt.Sprintf("redis://%s", _cluster.Addr()),
			cluster: false,
			want:    false,
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress(test.address),
			WithCluster(test.cluster),
		)

		if err != nil {
			t.Errorf("WithCluster returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Cluster, test.want) {
			t.Errorf("WithCluster is %v, want %v", _service.config.Cluster, test.want)
		}
	}
}
