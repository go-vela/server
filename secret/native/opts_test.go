// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/sqlite"
)

func TestNative_ClientOpt_WithDatabase(t *testing.T) {
	// setup types
	db, err := sqlite.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}
	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		database database.Service
		want     database.Service
	}{
		{
			failure:  false,
			database: db,
			want:     db,
		},
		{
			failure:  true,
			database: nil,
			want:     nil,
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithDatabase(test.database),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithDatabase should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithDatabase returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.Database, test.want) {
			t.Errorf("WithDatabase is %v, want %v", _service.Database, test.want)
		}
	}
}
