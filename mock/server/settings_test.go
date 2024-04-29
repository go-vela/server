// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/server/api/types/settings"
)

func TestSettings_GetResp(t *testing.T) {
	testSettings := settings.Platform{}

	err := json.Unmarshal([]byte(SettingsResp), &testSettings)
	if err != nil {
		t.Errorf("error unmarshaling settings: %v", err)
	}

	tSettings := reflect.TypeOf(testSettings)

	for i := 0; i < tSettings.NumField(); i++ {
		f := reflect.ValueOf(testSettings).Field(i)
		if f.IsNil() {
			t.Errorf("SettingsResp missing field %s", tSettings.Field(i).Name)
		}
	}
}

func TestSettings_UpdateResp(t *testing.T) {
	testSettings := settings.Platform{}

	err := json.Unmarshal([]byte(UpdateSettingsResp), &testSettings)
	if err != nil {
		t.Errorf("error unmarshaling settings: %v", err)
	}

	tSettings := reflect.TypeOf(testSettings)

	for i := 0; i < tSettings.NumField(); i++ {
		f := reflect.ValueOf(testSettings).Field(i)
		if f.IsNil() {
			t.Errorf("UpdateSettingsResp missing field %s", tSettings.Field(i).Name)
		}
	}
}
