// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_createHookService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  []string{CreateHookTable},
		Indexes: []string{CreateHookRepoIDNumberIndex, CreateHookRepoIDIndex},
	}

	// run test
	got := createHookService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createHookService is %v, want %v", got, want)
	}
}
