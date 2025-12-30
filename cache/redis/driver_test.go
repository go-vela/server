// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"github.com/go-vela/server/constants"
)

func TestRedis_Driver(t *testing.T) {
	// setup types
	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_redis, err := miniredis.Run()
	if err != nil {
		t.Errorf("unable to create miniredis instance: %v", err)
	}
	defer _redis.Close()

	want := constants.DriverRedis

	_service, err := New(
		t.Context(),
		WithAddress(fmt.Sprintf("redis://%s", _redis.Addr())),
		WithInstallTokenKey("123abc"),
		WithCluster(false),
	)
	if err != nil {
		t.Errorf("unable to create queue service: %v", err)
	}

	// run test
	got := _service.Driver()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Driver is %v, want %v", got, want)
	}
}
