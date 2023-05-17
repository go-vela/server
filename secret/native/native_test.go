// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/sqlite"
)

func TestNative_New(t *testing.T) {
	// setup types
	db, err := sqlite.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}

	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		database database.Interface
		want     database.Interface
	}{
		{
			failure:  false,
			database: db,
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
