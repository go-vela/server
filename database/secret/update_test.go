// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestSecret_Engine_UpdateSecret(t *testing.T) {
	// setup types
	_secretRepo := testSecret()
	_secretRepo.SetID(1)
	_secretRepo.SetOrg("foo")
	_secretRepo.SetRepo("bar")
	_secretRepo.SetName("baz")
	_secretRepo.SetValue("foob")
	_secretRepo.SetType("repo")
	_secretRepo.SetCreatedAt(1)
	_secretRepo.SetCreatedBy("user")
	_secretRepo.SetUpdatedAt(1)
	_secretRepo.SetUpdatedBy("user2")

	_secretOrg := testSecret()
	_secretOrg.SetID(2)
	_secretOrg.SetOrg("foo")
	_secretOrg.SetRepo("*")
	_secretOrg.SetName("bar")
	_secretOrg.SetValue("baz")
	_secretOrg.SetType("org")
	_secretOrg.SetCreatedAt(1)
	_secretOrg.SetCreatedBy("user")
	_secretOrg.SetUpdatedAt(1)
	_secretOrg.SetUpdatedBy("user2")

	_secretShared := testSecret()
	_secretShared.SetID(3)
	_secretShared.SetOrg("foo")
	_secretShared.SetTeam("bar")
	_secretShared.SetName("baz")
	_secretShared.SetValue("foob")
	_secretShared.SetType("shared")
	_secretShared.SetCreatedAt(1)
	_secretShared.SetCreatedBy("user")
	_secretShared.SetUpdatedAt(1)
	_secretShared.SetUpdatedBy("user2")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the repo query
	_mock.ExpectExec(`UPDATE "secrets"
SET "org"=$1,"repo"=$2,"team"=$3,"name"=$4,"value"=$5,"type"=$6,"images"=$7,"events"=$8,"allow_command"=$9,"created_at"=$10,"created_by"=$11,"updated_at"=$12,"updated_by"=$13
WHERE "id" = $14`).
		WithArgs("foo", "bar", nil, "baz", AnyArgument{}, "repo", nil, nil, false, 1, "user", AnyArgument{}, "user2", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// ensure the mock expects the org query
	_mock.ExpectExec(`UPDATE "secrets"
SET "org"=$1,"repo"=$2,"team"=$3,"name"=$4,"value"=$5,"type"=$6,"images"=$7,"events"=$8,"allow_command"=$9,"created_at"=$10,"created_by"=$11,"updated_at"=$12,"updated_by"=$13
WHERE "id" = $14`).
		WithArgs("foo", "*", nil, "bar", AnyArgument{}, "org", nil, nil, false, 1, "user", AnyArgument{}, "user2", 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// ensure the mock expects the shared query
	_mock.ExpectExec(`UPDATE "secrets"
SET "org"=$1,"repo"=$2,"team"=$3,"name"=$4,"value"=$5,"type"=$6,"images"=$7,"events"=$8,"allow_command"=$9,"created_at"=$10,"created_by"=$11,"updated_at"=$12,"updated_by"=$13
WHERE "id" = $14`).
		WithArgs("foo", nil, "bar", "baz", AnyArgument{}, "shared", nil, nil, false, 1, "user", NowTimestamp{}, "user2", 3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(context.TODO(), _secretRepo)
	if err != nil {
		t.Errorf("unable to create test repo secret for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSecret(context.TODO(), _secretOrg)
	if err != nil {
		t.Errorf("unable to create test org secret for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSecret(context.TODO(), _secretShared)
	if err != nil {
		t.Errorf("unable to create test shared secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		secret   *library.Secret
	}{
		{
			failure:  false,
			name:     "postgres with repo",
			database: _postgres,
			secret:   _secretRepo,
		},
		{
			failure:  false,
			name:     "postgres with org",
			database: _postgres,
			secret:   _secretOrg,
		},
		{
			failure:  false,
			name:     "postgres with shared",
			database: _postgres,
			secret:   _secretShared,
		},
		{
			failure:  false,
			name:     "sqlite3 with repo",
			database: _sqlite,
			secret:   _secretRepo,
		},
		{
			failure:  false,
			name:     "sqlite3 with org",
			database: _sqlite,
			secret:   _secretOrg,
		},
		{
			failure:  false,
			name:     "sqlite3 with shared",
			database: _sqlite,
			secret:   _secretShared,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.UpdateSecret(context.TODO(), test.secret)
			got.SetUpdatedAt(test.secret.GetUpdatedAt())

			if test.failure {
				if err == nil {
					t.Errorf("UpdateSecret for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateSecret for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.secret) {
				t.Errorf("UpdateSecret for %s is %s, want %s", test.name, got, test.secret)
			}
		})
	}
}
