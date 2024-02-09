// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestRepo_ActiveRepoResp(t *testing.T) {
	testRepo := library.Repo{}

	err := json.Unmarshal([]byte(RepoResp), &testRepo)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tRepo := reflect.TypeOf(testRepo)

	for i := 0; i < tRepo.NumField(); i++ {
		if tRepo.Field(i).Name == "Hash" {
			continue
		}
		if reflect.ValueOf(testRepo).Field(i).IsNil() {
			t.Errorf("RepoResp missing field %s", tRepo.Field(i).Name)
		}
	}
}
