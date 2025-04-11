// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestDeployment_ActiveDeploymentResp(t *testing.T) {
	testDeployment := api.Deployment{}

	err := json.Unmarshal([]byte(DeploymentResp), &testDeployment)
	if err != nil {
		t.Errorf("error unmarshaling deployment: %v", err)
	}

	tDeployment := reflect.TypeOf(testDeployment)

	for i := range tDeployment.NumField() {
		if reflect.ValueOf(testDeployment).Field(i).IsNil() {
			t.Errorf("DeploymentResp missing field %s", tDeployment.Field(i).Name)
		}
	}
}
