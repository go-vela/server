// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestHook_ActiveHookResp(t *testing.T) {
	testHook := api.Hook{}

	err := json.Unmarshal([]byte(HookResp), &testHook)
	if err != nil {
		t.Errorf("error unmarshaling hook: %v", err)
	}

	tHook := reflect.TypeOf(testHook)

	for i := 0; i < tHook.NumField(); i++ {
		if reflect.ValueOf(testHook).Field(i).IsNil() {
			t.Errorf("HookResp missing field %s", tHook.Field(i).Name)
		}
	}
}
