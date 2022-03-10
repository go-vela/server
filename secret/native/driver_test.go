// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/types/constants"
)

func TestNative_Driver(t *testing.T) {
	// setup types
	db, err := sqlite.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}

	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	want := constants.DriverNative

	_service, err := New(
		WithDatabase(db),
	)
	if err != nil {
		t.Errorf("unable to create secret service: %v", err)
	}

	// run test
	got := _service.Driver()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Driver is %v, want %v", got, want)
	}
}
