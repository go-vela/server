// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestStep_ActiveStepResp(t *testing.T) {
	testStep := library.Step{}

	err := json.Unmarshal([]byte(StepResp), &testStep)
	if err != nil {
		t.Errorf("error unmarshaling step: %v", err)
	}

	tStep := reflect.TypeOf(testStep)

	for i := 0; i < tStep.NumField(); i++ {
		if reflect.ValueOf(testStep).Field(i).IsNil() {
			t.Errorf("StepResp missing field %s", tStep.Field(i).Name)
		}
	}
}
