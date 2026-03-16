// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestRepo_Engine_PartialUpdateRepo(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.GetOwner().SetID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetPreviousName("oldName")
	_repo.SetApproveBuild(constants.ApproveForkAlways)
	_repo.SetTopics([]string{})
	_repo.SetAllowEvents(api.NewEventsFromMask(1))
	_repo.SetApprovalTimeout(5)
	_repo.SetCustomProps(map[string]any{"foo": "bar"})

	updates := new(api.Repo)
	updates.SetID(1)
	updates.SetTopics([]string{"topic1", "topic2"})
	updates.SetVisibility("private")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "repos" SET "topics"=$1,"visibility"=$2 WHERE "id" = $3`).
		WithArgs(`{"topic1","topic2"}`, "private", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateRepo(context.TODO(), _repo)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
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
			err := test.database.PartialUpdateRepo(context.TODO(), updates)

			if test.failure {
				if err == nil {
					t.Errorf("PartialUpdateRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("PartialUpdateRepo for %s returned err: %v", test.name, err)
			}
		})
	}
}
