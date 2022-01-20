// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestSqlite_Client_Driver(t *testing.T) {
	// setup types
	want := constants.DriverSqlite

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// run test
	got := _database.Driver()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Driver is %v, want %v", got, want)
	}
}
