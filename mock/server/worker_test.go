// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/server/api/types"
)

func TestWorker_ActiveWorkerResp(t *testing.T) {
	testWorker := types.Worker{}

	err := json.Unmarshal([]byte(WorkerResp), &testWorker)
	if err != nil {
		t.Errorf("error unmarshaling worker: %v", err)
	}

	tWorker := reflect.TypeOf(testWorker)

	for i := 0; i < tWorker.NumField(); i++ {
		if reflect.ValueOf(testWorker).Field(i).IsNil() {
			t.Errorf("WorkerResp missing field %s", tWorker.Field(i).Name)
		}
	}
}

func TestWorker_ListActiveWorkerResp(t *testing.T) {
	testWorkers := []types.Worker{}

	err := json.Unmarshal([]byte(WorkersResp), &testWorkers)
	if err != nil {
		t.Errorf("error unmarshaling worker: %v", err)
	}

	for index, worker := range testWorkers {
		tWorker := reflect.TypeOf(worker)

		for i := 0; i < tWorker.NumField(); i++ {
			if reflect.ValueOf(worker).Field(i).IsNil() {
				t.Errorf("WorkersResp index %d missing field %s", index, tWorker.Field(i).Name)
			}
		}
	}
}
