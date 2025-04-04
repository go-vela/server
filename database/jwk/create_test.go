// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestJWK_Engine_CreateJWK(t *testing.T) {
	// setup types
	_jwk := testutils.JWK()
	_jwkBytes, err := json.Marshal(_jwk)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	kid, ok := _jwk.KeyID()
	if !ok {
		t.Errorf("unable to get key ID for jwk")
	}

	// ensure the mock expects the query
	_mock.ExpectExec(`INSERT INTO "jwks"
("id","active","key")
VALUES ($1,$2,$3)`).
		WithArgs(kid, true, _jwkBytes).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			err := test.database.CreateJWK(context.TODO(), _jwk)

			if test.failure {
				if err == nil {
					t.Errorf("CreateDashboard for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateDashboard for %s returned err: %v", test.name, err)
			}
		})
	}
}
