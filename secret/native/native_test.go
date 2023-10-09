// SPDX-License-Identifier: Apache-2.0

package native

import (
	"testing"

	"github.com/go-vela/server/database"
)

func TestNative_New(t *testing.T) {
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
		},
		{
			failure:  true,
			database: nil,
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			WithDatabase(test.database),
		)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}
