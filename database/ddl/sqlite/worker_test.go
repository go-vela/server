// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_createWorkerService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  []string{CreateWorkerTable},
		Indexes: []string{CreateWorkersHostnameAddressIndex},
	}

	// run test
	got := createWorkerService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createWorkerervice is %v, want %v", got, want)
	}
}
