// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/database/testutils"
)

func TestRepo_Engine_CreateRepo(t *testing.T) {
	// setup types
	props := map[string]any{
		"foo": "bar",
	}

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
	_repo.SetTopics([]string{})
	_repo.SetCustomProps(props)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "repos"
("user_id","hash","org","name","full_name","link","clone","branch","topics","build_limit","timeout","counter","visibility","private","trusted","active","allow_events","pipeline_type","previous_name","approve_build","approval_timeout","install_id","custom_props","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24) RETURNING "id"`).
		WithArgs(1, AnyArgument{}, "foo", "bar", "foo/bar", nil, nil, nil, AnyArgument{}, AnyArgument{}, AnyArgument{}, AnyArgument{}, "public", false, false, false, nil, "yaml", "oldName", nil, nil, 0, `{"foo":"bar"}`, 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			got, err := test.database.CreateRepo(context.TODO(), _repo)

			if test.failure {
				if err == nil {
					t.Errorf("CreateRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateRepo for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(_repo, got); diff != "" {
				t.Errorf("CreateRepo mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
