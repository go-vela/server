// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestSchedule_ActiveScheduleResp(t *testing.T) {
	testSchedule := api.Schedule{}

	err := json.Unmarshal([]byte(ScheduleResp), &testSchedule)
	if err != nil {
		t.Errorf("error unmarshaling schedule: %v", err)
	}

	tSchedule := reflect.TypeOf(testSchedule)

	for i := range tSchedule.NumField() {
		if reflect.ValueOf(testSchedule).Field(i).IsNil() {
			t.Errorf("ScheduleResp missing field %s", tSchedule.Field(i).Name)
		}
	}
}
