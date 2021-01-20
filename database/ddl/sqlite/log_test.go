// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_createLogService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  CreateLogTable,
		Indexes: []string{CreateLogBuildIDIndex, CreateLogStepIDIndex, CreateLogServiceIDIndex},
	}

	// run test
	got := createLogService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createLogService is %v, want %v", got, want)
	}
}
