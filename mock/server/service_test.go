// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestService_ActiveServiceResp(t *testing.T) {
	testService := library.Service{}

	err := json.Unmarshal([]byte(ServiceResp), &testService)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tService := reflect.TypeOf(testService)

	for i := 0; i < tService.NumField(); i++ {
		if reflect.ValueOf(testService).Field(i).IsNil() {
			t.Errorf("ServiceResp missing field %s", tService.Field(i).Name)
		}
	}
}
