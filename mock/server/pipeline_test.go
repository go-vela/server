// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestPipeline_ActivePipelineResp(t *testing.T) {
	testPipeline := api.Pipeline{}

	err := json.Unmarshal([]byte(PipelineResp), &testPipeline)
	if err != nil {
		t.Errorf("error unmarshaling pipeline: %v", err)
	}

	tPipeline := reflect.TypeFor[api.Pipeline]()

	for i := 0; i < tPipeline.NumField(); i++ {
		if reflect.ValueOf(testPipeline).Field(i).IsNil() {
			t.Errorf("PipelineResp missing field %s", tPipeline.Field(i).Name)
		}
	}
}
