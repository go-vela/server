// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
		List: map[string]string{
			"all":    ListSecrets,
			"org":    ListOrgSecrets,
			"repo":   ListRepoSecrets,
			"shared": ListSharedSecrets,
		},
		Select: map[string]string{
			"org":         SelectOrgSecret,
			"repo":        SelectRepoSecret,
			"shared":      SelectSharedSecret,
			"countOrg":    SelectOrgSecretsCount,
			"countRepo":   SelectRepoSecretsCount,
			"countShared": SelectSharedSecretsCount,
		},
		Delete: DeleteSecret,
	}

	// run test
	got := createSecretService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createSecretService is %v, want %v", got, want)
	}
}
