// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestWorker_ActiveWorkerResp(t *testing.T) {
	testWorker := library.Worker{}

	err := json.Unmarshal([]byte(WorkerResp), &testWorker)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tWorker := reflect.TypeOf(testWorker)

	for i := 0; i < tWorker.NumField(); i++ {
		if reflect.ValueOf(testWorker).Field(i).IsNil() {
			t.Errorf("WorkerResp missing field %s", tWorker.Field(i).Name)
		}
	}
}
