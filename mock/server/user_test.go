// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestUser_ActiveUserResp(t *testing.T) {
	testUser := api.User{}

	err := json.Unmarshal([]byte(UserResp), &testUser)
	if err != nil {
		t.Errorf("error unmarshaling user: %v", err)
	}

	tUser := reflect.TypeOf(testUser)

	for i := 0; i < tUser.NumField(); i++ {
		if tUser.Field(i).Name == "Token" || tUser.Field(i).Name == "RefreshToken" || tUser.Field(i).Name == "Hash" {
			continue
		}

		if reflect.ValueOf(testUser).Field(i).IsNil() {
			t.Errorf("UserResp missing field %s", tUser.Field(i).Name)
		}
	}
}
