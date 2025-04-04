// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestJWK_Engine_RotateKeys(t *testing.T) {
	// setup types
	_jwkOne := testutils.JWK()
	_jwkOneBytes, err := json.Marshal(_jwkOne)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	_jwkTwo := testutils.JWK()
	_jwkTwoBytes, err := json.Marshal(_jwkTwo)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	kidOne, ok := _jwkOne.KeyID()
	if !ok {
		t.Errorf("unable to get key ID for jwk")
	}

	kidTwo, ok := _jwkTwo.KeyID()
	if !ok {
		t.Errorf("unable to get key ID for jwk")
	}

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "active", "key"},
	).AddRow(kidOne, true, _jwkOneBytes)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks" WHERE id = $1 AND active = $2 LIMIT $3`).WithArgs(kidOne, true, 1).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "active", "key"},
	).AddRow(kidTwo, true, _jwkTwoBytes)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks" WHERE id = $1 AND active = $2 LIMIT $3`).WithArgs(kidTwo, true, 1).WillReturnRows(_rows)

	_mock.ExpectExec(`DELETE FROM "jwks" WHERE active = $1`).
		WithArgs(false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`UPDATE "jwks" SET "active"=$1 WHERE active = $2`).
		WithArgs(false, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err = _sqlite.CreateJWK(context.TODO(), _jwkOne)
	if err != nil {
		t.Errorf("unable to create test jwk for sqlite: %v", err)
	}

	err = _sqlite.CreateJWK(context.TODO(), _jwkTwo)
	if err != nil {
		t.Errorf("unable to create test jwk for sqlite: %v", err)
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
			_, err := test.database.GetActiveJWK(context.TODO(), kidOne)
			if err != nil {
				t.Errorf("GetActiveJWK for %s returned err: %v", test.name, err)
			}

			_, err = test.database.GetActiveJWK(context.TODO(), kidTwo)
			if err != nil {
				t.Errorf("GetActiveJWK for %s returned err: %v", test.name, err)
			}

			err = test.database.RotateKeys(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("RotateKeys for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("RotateKeys for %s returned err: %v", test.name, err)
			}

			_, err = test.database.GetActiveJWK(context.TODO(), kidOne)
			if err == nil {
				t.Errorf("GetActiveJWK for %s should have returned err", test.name)
			}

			_, err = test.database.GetActiveJWK(context.TODO(), kidTwo)
			if err == nil {
				t.Errorf("GetActiveJWK for %s should have returned err", test.name)
			}
		})
	}
}
