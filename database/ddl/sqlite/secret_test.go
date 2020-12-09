// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_createSecretService(t *testing.T) {
	// setup types
	want := &Service{
		Create: CreateSecretTable,
		Indexes: []string{
			CreateSecretTypeOrgRepo,
			CreateSecretTypeOrgTeam,
			CreateSecretTypeOrg,
			CreateSecretType,
		},
	}

	// run test
	got := createSecretService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createSecretService is %v, want %v", got, want)
	}
}
