// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestGitHub_Driver(t *testing.T) {
	// setup types
	want := constants.DriverGithub

	_service, err := New(
		WithAddress("https://github.com/"),
		WithClientID("foo"),
		WithClientSecret("bar"),
		WithServerAddress("https://vela-server.example.com"),
		WithStatusContext("continuous-integration/vela"),
		WithWebUIAddress("https://vela.example.com"),
	)
	if err != nil {
		t.Errorf("unable to create scm service: %v", err)
	}

	// run test
	got := _service.Driver()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Driver is %v, want %v", got, want)
	}
}
