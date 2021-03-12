// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package source

import (
	"reflect"
	"testing"
)

func TestSource_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
		want    Service
	}{
		{
			failure: true,
			setup: &Setup{
				Driver:        "github",
				Address:       "https://github.com",
				ClientID:      "foo",
				ClientSecret:  "bar",
				ServerAddress: "https://vela-server.example.com",
				StatusContext: "continuous-integration/vela",
				WebUIAddress:  "https://vela.example.com",
			},
			want: nil,
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "gitlab",
				Address:       "https://gitlab.com",
				ClientID:      "foo",
				ClientSecret:  "bar",
				ServerAddress: "https://vela-server.example.com",
				StatusContext: "continuous-integration/vela",
				WebUIAddress:  "https://vela.example.com",
			},
			want: nil,
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "bitbucket",
				Address:       "https://bitbucket.org",
				ClientID:      "foo",
				ClientSecret:  "bar",
				ServerAddress: "https://vela-server.example.com",
				StatusContext: "continuous-integration/vela",
				WebUIAddress:  "https://vela.example.com",
			},
			want: nil,
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "github",
				Address:       "",
				ClientID:      "foo",
				ClientSecret:  "bar",
				ServerAddress: "https://vela-server.example.com",
				StatusContext: "continuous-integration/vela",
				WebUIAddress:  "https://vela.example.com",
			},
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := New(test.setup)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("New is %v, want %v", got, test.want)
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("New is %v, want %v", got, test.want)
		}
	}
}
