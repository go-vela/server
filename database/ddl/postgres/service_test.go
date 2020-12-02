// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createServiceService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  []string{CreateServiceTable},
		Indexes: []string{CreateServiceBuildIDNumberIndex},
	}

	// run test
	got := createServiceService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createServiceService is %v, want %v", got, want)
	}
}
