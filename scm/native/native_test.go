// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"testing"
)

func TestNative_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		id      string
	}{
		{
			failure: false,
			id:      "foo",
		},
		{
			failure: true,
			id:      "",
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			WithAddress("https://github.com/"),
			WithClientID(test.id),
			WithClientSecret("bar"),
			WithKind("github"),
			WithServerAddress("https://vela-server.example.com"),
			WithStatusContext("continuous-integration/vela"),
			WithWebUIAddress("https://vela.example.com"),
			WithScopes([]string{"repo", "repo:status", "user:email", "read:user", "read:org"}),
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
