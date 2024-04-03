// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestDeployment_ActiveDeploymentResp(t *testing.T) {
	testDeployment := library.Deployment{}

	err := json.Unmarshal([]byte(DeploymentResp), &testDeployment)
	if err != nil {
		t.Errorf("error unmarshaling deployment: %v", err)
	}

	tDeployment := reflect.TypeOf(testDeployment)

	for i := 0; i < tDeployment.NumField(); i++ {
		if reflect.ValueOf(testDeployment).Field(i).IsNil() {
			t.Errorf("DeploymentResp missing field %s", tDeployment.Field(i).Name)
		}
	}
}
