// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestHook_ActiveHookResp(t *testing.T) {
	testHook := library.Hook{}

	err := json.Unmarshal([]byte(HookResp), &testHook)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tHook := reflect.TypeOf(testHook)

	for i := 0; i < tHook.NumField(); i++ {
		if reflect.ValueOf(testHook).Field(i).IsNil() {
			t.Errorf("HookResp missing field %s", tHook.Field(i).Name)
		}
	}
}
