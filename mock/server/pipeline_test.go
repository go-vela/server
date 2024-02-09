// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestPipeline_ActivePipelineResp(t *testing.T) {
	testPipeline := library.Pipeline{}

	err := json.Unmarshal([]byte(PipelineResp), &testPipeline)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tPipeline := reflect.TypeOf(testPipeline)

	for i := 0; i < tPipeline.NumField(); i++ {
		if reflect.ValueOf(testPipeline).Field(i).IsNil() {
			t.Errorf("PipelineResp missing field %s", tPipeline.Field(i).Name)
		}
	}
}
