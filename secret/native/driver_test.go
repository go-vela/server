// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
)

func TestNative_Driver(t *testing.T) {
	// setup types
	d, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}
	defer d.Database.Close()

	want := constants.DriverNative

	_service, err := New(
		WithDatabase(d),
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
