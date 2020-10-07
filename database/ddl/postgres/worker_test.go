// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createWorkerService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  CreateWorkerTable,
		Indexes: []string{CreateWorkersHostnameAddressIndex},
	}

	// run test
	got := createWorkerService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createWorkerService is %v, want %v", got, want)
	}
}
