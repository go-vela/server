// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createStepService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  CreateStepTable,
		Indexes: []string{CreateStepBuildIDNumberIndex},
	}

	// run test
	got := createStepService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createStepService is %v, want %v", got, want)
	}
}
