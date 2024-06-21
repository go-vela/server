// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/go-vela/server/database/testutils"
)

func TestTypes_JWK_Nullify(t *testing.T) {
	// setup tests
	var j *JWK

	tests := []struct {
		JWK  *JWK
		want *JWK
	}{
		{
			JWK:  testJWK(),
			want: testJWK(),
		},
		{
			JWK:  j,
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.JWK.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_JWK_ToAPI(t *testing.T) {
	// setup types
	want := testutils.JWK()

	wantBytes, err := json.Marshal(want)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	uuid, _ := uuid.Parse(want.KeyID())
	h := &JWK{
		ID:     uuid,
		Active: sql.NullBool{Bool: true, Valid: true},
		Key:    wantBytes,
	}

	// run test
	got := h.ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_JWKFromAPI(t *testing.T) {
	j := testutils.JWK()

	jBytes, err := json.Marshal(j)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	uuid, err := uuid.Parse(j.KeyID())
	if err != nil {
		t.Errorf("unable to parse JWK key id: %v", err)
	}

	// setup types
	want := &JWK{
		ID:     uuid,
		Active: sql.NullBool{Bool: false, Valid: false},
		Key:    jBytes,
	}

	// run test
	got := JWKFromAPI(j)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("JWKFromAPI() mismatch (-want +got):\n%s", diff)
	}
}

// testJWK is a test helper function to create a JWK
// type with all fields set to a fake value.
func testJWK() *JWK {
	uuid, _ := uuid.Parse("c8da1302-07d6-11ea-882f-4893bca275b8")

	return &JWK{
		ID:     uuid,
		Active: sql.NullBool{Bool: true, Valid: true},
	}
}
