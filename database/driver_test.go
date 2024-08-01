// SPDX-License-Identifier: Apache-2.0

package database

import (
	"strings"
	"testing"
)

func TestDatabase_Engine_Driver(t *testing.T) {
	_postgres, _ := testPostgres(t)
	defer _postgres.Close()

	_sqlite := testSqlite(t)
	defer _sqlite.Close()

	// setup tests
	tests := []struct {
		name     string
		database *engine
		want     string
	}{
		{
			name:     "success with postgres",
			database: _postgres,
			want:     "postgres",
		},
		{
			name:     "success with sqlite3",
			database: _sqlite,
			want:     "sqlite3",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.database.Driver()

			if !strings.EqualFold(got, test.want) {
				t.Errorf("Driver for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
