// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestDashboard_ActiveDashboardResp(t *testing.T) {
	testDashboard := api.Dashboard{}

	err := json.Unmarshal([]byte(DashboardResp), &testDashboard)
	if err != nil {
		t.Errorf("error unmarshaling build: %v", err)
	}

	tDashboard := reflect.TypeOf(testDashboard)

	for i := 0; i < tDashboard.NumField(); i++ {
		if reflect.ValueOf(testDashboard).Field(i).IsNil() {
			t.Errorf("DashboardResp missing field %s", tDashboard.Field(i).Name)
		}
	}

	testDashboards := []api.Dashboard{}
	err = json.Unmarshal([]byte(DashboardsResp), &testDashboards)
	if err != nil {
		t.Errorf("error unmarshaling builds: %v", err)
	}

	for _, testDashboard := range testDashboards {
		tDashboard := reflect.TypeOf(testDashboard)

		for i := 0; i < tDashboard.NumField(); i++ {
			if reflect.ValueOf(testDashboard).Field(i).IsNil() {
				t.Errorf("DashboardsResp missing field %s", tDashboard.Field(i).Name)
			}
		}
	}
}
