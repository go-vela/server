// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
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
			context.Background(),
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

func TestRedis_ClientOpt_WithInstallTokenKey(t *testing.T) {
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
		failure    bool
		installKey string
		want       string
	}{
		{
			failure:    false,
			installKey: "foo",
			want:       "foo",
		},
		{
			failure:    true,
			installKey: "",
			want:       "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			t.Context(),
			WithAddress(fmt.Sprintf("redis://%s", _redis.Addr())),
			WithInstallTokenKey(test.installKey),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithInstallTokenKey should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithInstallTokenKey returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.InstallTokenKey, test.want) {
			t.Errorf("WithInstallTokenKey is %v, want %v", _service.config.InstallTokenKey, test.want)
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
			t.Context(),
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
