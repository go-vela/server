// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSecret_Engine_UpdateSecret(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy("user")
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy("user2")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "secrets"
SET "org"=$1,"repo"=$2,"team"=$3,"name"=$4,"value"=$5,"type"=$6,"images"=$7,"events"=$8,"allow_command"=$9,"created_at"=$10,"created_by"=$11,"updated_at"=$12,"updated_by"=$13
WHERE "id" = $14`).
		WithArgs("foo", "bar", nil, "baz", AnyArgument{}, "repo", nil, nil, false, 1, "user", time.Now().UTC().Unix(), "user2", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateSecret(_secret)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = test.database.UpdateSecret(_secret)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateSecret for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateSecret for %s returned err: %v", test.name, err)
			}
		})
	}
}
