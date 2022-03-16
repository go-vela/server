// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestVault_Driver(t *testing.T) {
	// setup types
	want := constants.DriverVault

	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	type args struct {
		version string
		prefix  string
	}

	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_service, err := New(
				WithAddress(fake.URL),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
			if err != nil {
				t.Errorf("unable to create secret service: %v", err)
			}

			got := _service.Driver()

			if !reflect.DeepEqual(got, want) {
				t.Errorf("Driver is %v, want %v", got, want)
			}
		})
	}
}
