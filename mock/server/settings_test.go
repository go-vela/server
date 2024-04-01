// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestSettings_GetResp(t *testing.T) {
	testSettings := api.Settings{}

	err := json.Unmarshal([]byte(SettingsResp), &testSettings)
	if err != nil {
		t.Errorf("error unmarshaling settings: %v", err)
	}

	tSettings := reflect.TypeOf(testSettings)

	for i := 0; i < tSettings.NumField(); i++ {
		if reflect.ValueOf(testSettings).Field(i).IsNil() {
			t.Errorf("SettingsResp missing field %s", tSettings.Field(i).Name)
		}
	}
}
