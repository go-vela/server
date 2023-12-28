// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBuild_Engine_CreateBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "builds"
("repo_id","pipeline_id","number","parent","event","event_action","status","error","enqueued","created","started","finished","deploy","deploy_payload","clone","source","title","message","commit","sender","author","email","link","branch","ref","base_ref","head_ref","host","runtime","distribution","approved_at","approved_by","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33) RETURNING "id"`).
		WithArgs(1, nil, 1, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, AnyArgument{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			got, err := test.database.CreateBuild(context.TODO(), _build)

			if test.failure {
				if err == nil {
					t.Errorf("CreateBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _build) {
				t.Errorf("CreateBuild for %s returned %s, want %s", test.name, got, _build)
			}
		})
	}
}
