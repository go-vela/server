// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSecret_Engine_CountSecretsForOrg(t *testing.T) {
	// setup types
	_secretOne := TestSecret()
	_secretOne.SetID(1)
	_secretOne.SetUserID(1)
	_secretOne.SetHash("baz")
	_secretOne.SetOrg("foo")
	_secretOne.SetName("bar")
	_secretOne.SetFullName("foo/bar")
	_secretOne.SetVisibility("public")

	_secretTwo := TestSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetUserID(1)
	_secretTwo.SetHash("baz")
	_secretTwo.SetOrg("bar")
	_secretTwo.SetName("foo")
	_secretTwo.SetFullName("bar/foo")
	_secretTwo.SetVisibility("public")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "secrets" WHERE org = $1`).WithArgs("foo").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateSecret(_secretOne)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	err = _sqlite.CreateSecret(_secretTwo)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     1,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     1,
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountSecretsForOrg("foo", filters)

			if test.failure {
				if err == nil {
					t.Errorf("CountSecretsForOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountSecretsForOrg for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountSecretsForOrg for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}