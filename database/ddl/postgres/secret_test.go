// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createSecretService(t *testing.T) {
	// setup types
	want := &Service{
		Create: []string{CreateSecretTable},
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
