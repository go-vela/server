// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestRepo_ActiveRepoResp(t *testing.T) {
	testRepo := api.Repo{}

	err := json.Unmarshal([]byte(RepoResp), &testRepo)
	if err != nil {
		t.Errorf("error unmarshaling repo: %v", err)
	}

	tRepo := reflect.TypeFor[api.Repo]()

	for i := 0; i < tRepo.NumField(); i++ {
		if tRepo.Field(i).Name == "Hash" {
			continue
		}

		if reflect.ValueOf(testRepo).Field(i).IsNil() {
			t.Errorf("RepoResp missing field %s", tRepo.Field(i).Name)
		}
	}
}
