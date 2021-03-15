// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"testing"

	"github.com/go-vela/server/database"
)

func TestNative_New(t *testing.T) {
	// setup types
	d, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}
	defer d.Database.Close()

	// setup tests
	tests := []struct {
		failure  bool
		database database.Service
		want     database.Service
	}{
		{
			failure:  false,
			database: d,
		},
		{
			failure:  true,
			database: nil,
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			WithDatabase(test.database),
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
