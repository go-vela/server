// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestSecret_ActiveSecretResp(t *testing.T) {
	testSecret := api.Secret{}

	err := json.Unmarshal([]byte(SecretResp), &testSecret)
	if err != nil {
		t.Errorf("error unmarshaling secret: %v", err)
	}

	tSecret := reflect.TypeFor[api.Secret]()

	for i := 0; i < tSecret.NumField(); i++ {
		if reflect.ValueOf(testSecret).Field(i).IsNil() {
			t.Errorf("SecretResp missing field %s", tSecret.Field(i).Name)
		}
	}
}
