// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestLog_ActiveLogResp(t *testing.T) {
	testLog := library.Log{}

	err := json.Unmarshal([]byte(LogResp), &testLog)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tLog := reflect.TypeOf(testLog)

	for i := 0; i < tLog.NumField(); i++ {
		if reflect.ValueOf(testLog).Field(i).IsNil() {
			t.Errorf("LogResp missing field %s", tLog.Field(i).Name)
		}
	}
}
