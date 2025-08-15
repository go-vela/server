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
		t.Errorf("error unmarshaling dashboard: %v", err)
	}

	tDashboard := reflect.TypeOf(testDashboard)

	for i := 0; i < tDashboard.NumField(); i++ {
		if reflect.ValueOf(testDashboard).Field(i).IsNil() {
			t.Errorf("DashboardResp missing field %s", tDashboard.Field(i).Name)
		}
	}

	testDashCard := api.DashCard{}

	err = json.Unmarshal([]byte(DashCardResp), &testDashCard)
	if err != nil {
		t.Errorf("error unmarshaling dash card: %v", err)
	}

	tDashCard := reflect.TypeOf(testDashCard)

	for i := 0; i < tDashCard.NumField(); i++ {
		if reflect.ValueOf(testDashCard).Field(i).IsNil() {
			t.Errorf("DashCardResp missing field %s", tDashCard.Field(i).Name)
		}
	}

	testDashCards := []api.DashCard{}

	err = json.Unmarshal([]byte(DashCardsResp), &testDashCards)
	if err != nil {
		t.Errorf("error unmarshaling dash cards: %v", err)
	}

	for _, testDashCard := range testDashCards {
		tDashCard := reflect.TypeOf(testDashCard)

		for i := 0; i < tDashCard.NumField(); i++ {
			if reflect.ValueOf(testDashCard).Field(i).IsNil() {
				t.Errorf("DashboardsResp missing field %s", tDashboard.Field(i).Name)
			}
		}
	}
}
