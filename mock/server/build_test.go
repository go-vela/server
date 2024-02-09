// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestBuild_ActiveBuildResp(t *testing.T) {
	testBuild := library.Build{}

	err := json.Unmarshal([]byte(BuildResp), &testBuild)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tBuild := reflect.TypeOf(testBuild)

	for i := 0; i < tBuild.NumField(); i++ {
		if reflect.ValueOf(testBuild).Field(i).IsNil() {
			t.Errorf("BuildResp missing field %s", tBuild.Field(i).Name)
		}
	}
}
