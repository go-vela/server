// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSchedule_ActiveScheduleResp(t *testing.T) {
	testSchedule := library.Schedule{}

	err := json.Unmarshal([]byte(ScheduleResp), &testSchedule)
	if err != nil {
		t.Errorf("error unmarshaling schedule: %v", err)
	}

	tSchedule := reflect.TypeOf(testSchedule)

	for i := 0; i < tSchedule.NumField(); i++ {
		if reflect.ValueOf(testSchedule).Field(i).IsNil() {
			t.Errorf("ScheduleResp missing field %s", tSchedule.Field(i).Name)
		}
	}
}
