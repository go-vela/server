// SPDX-License-Identifier: Apache-2.0

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
)

func TestNative_Driver(t *testing.T) {
	// setup types
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}
	defer db.Close()

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
