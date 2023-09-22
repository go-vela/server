// SPDX-License-Identifier: Apache-2.0

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
)

func TestNative_ClientOpt_WithDatabase(t *testing.T) {
	// setup types
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	// setup tests
	tests := []struct {
		failure  bool
		database database.Interface
		want     database.Interface
	}{
		{
			failure:  false,
			database: db,
			want:     db,
		},
		{
			failure:  true,
			database: nil,
			want:     nil,
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithDatabase(test.database),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithDatabase should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithDatabase returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.Database, test.want) {
			t.Errorf("WithDatabase is %v, want %v", _service.Database, test.want)
		}
	}
}
