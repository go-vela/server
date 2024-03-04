// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSecret_ActiveSecretResp(t *testing.T) {
	testSecret := library.Secret{}

	err := json.Unmarshal([]byte(SecretResp), &testSecret)
	if err != nil {
		t.Errorf("error unmarshaling secret: %v", err)
	}

	tSecret := reflect.TypeOf(testSecret)

	for i := 0; i < tSecret.NumField(); i++ {
		if reflect.ValueOf(testSecret).Field(i).IsNil() {
			t.Errorf("SecretResp missing field %s", tSecret.Field(i).Name)
		}
	}
}
