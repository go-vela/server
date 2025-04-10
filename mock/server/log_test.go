// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestLog_ActiveLogResp(t *testing.T) {
	testLog := api.Log{}

	err := json.Unmarshal([]byte(LogResp), &testLog)
	if err != nil {
		t.Errorf("error unmarshaling log: %v", err)
	}

	tLog := reflect.TypeOf(testLog)

	for i := range tLog.NumField() {
		if reflect.ValueOf(testLog).Field(i).IsNil() {
			t.Errorf("LogResp missing field %s", tLog.Field(i).Name)
		}
	}
}
