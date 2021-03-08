// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package source

import (
	"testing"
)

func TestSource_Setup_Github(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:        "github",
		Address:       "https://github.com",
		ClientID:      "foo",
		ClientSecret:  "bar",
		StatusContext: "continuous-integration/vela",
	}

	got, err := _setup.Github()
	if err == nil {
		t.Errorf("Github should have returned err")
	}

	if got != nil {
		t.Errorf("Github is %v, want nil", got)
	}
}

func TestSource_Setup_Gitlab(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:        "gitlab",
		Address:       "https://gitlab.com",
		ClientID:      "foo",
		ClientSecret:  "bar",
		StatusContext: "continuous-integration/vela",
	}

	got, err := _setup.Gitlab()
	if err == nil {
		t.Errorf("Gitlab should have returned err")
	}

	if got != nil {
		t.Errorf("Gitlab is %v, want nil", got)
	}
}

func TestSource_Setup_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Driver:        "github",
				Address:       "https://github.com",
				ClientID:      "foo",
				ClientSecret:  "bar",
				StatusContext: "continuous-integration/vela",
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:        "gitlab",
				Address:       "https://gitlab.com",
				ClientID:      "foo",
				ClientSecret:  "bar",
				StatusContext: "continuous-integration/vela",
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
