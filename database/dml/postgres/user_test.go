// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createUserService(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":  ListUsers,
			"lite": ListLiteUsers,
		},
		Select: map[string]string{
			"user":         SelectUser,
			"name":         SelectUserName,
			"count":        SelectUsersCount,
			"refreshToken": SelectRefreshToken,
		},
		Delete: DeleteUser,
	}

	// run test
	got := createUserService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createUserService is %v, want %v", got, want)
	}
}
