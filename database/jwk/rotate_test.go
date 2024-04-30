// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestJWK_Engine_RotateKeys(t *testing.T) {
	// setup types
	_jwkOne := testutils.APIJWK()
	_jwkOne.Kid = "c8da1302-07d6-11ea-882f-4893bca275b8"
	_jwkOne.Algorithm = "RS256"
	_jwkOne.Kty = "rsa"
	_jwkOne.Use = "sig"
	_jwkOne.N = "123456"
	_jwkOne.E = "123"

	_jwkTwo := testutils.APIJWK()
	_jwkTwo.Kid = "c8da1302-07d6-11ea-882f-4893bca275b9"
	_jwkTwo.Algorithm = "RS256"
	_jwkTwo.Kty = "rsa"
	_jwkTwo.Use = "sig"
	_jwkTwo.N = "123789"
	_jwkTwo.E = "456"

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "active", "key"},
	).AddRow("c8da1302-07d6-11ea-882f-4893bca275b8", true, []byte(`{"alg":"RS256","use":"sig","x5t":"","kid":"c8da1302-07d6-11ea-882f-4893bca275b8","kty":"rsa","x5c":null,"n":"123456","e":"123"}`))

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks" WHERE id = $1 AND active = $2 LIMIT $3`).WithArgs("c8da1302-07d6-11ea-882f-4893bca275b8", true, 1).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "active", "key"},
	).AddRow("c8da1302-07d6-11ea-882f-4893bca275b9", true, []byte(`{"alg":"RS256","use":"sig","x5t":"","kid":"c8da1302-07d6-11ea-882f-4893bca275b9","kty":"rsa","x5c":null,"n":"123789","e":"456"}`))

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks" WHERE id = $1 AND active = $2 LIMIT $3`).WithArgs("c8da1302-07d6-11ea-882f-4893bca275b9", true, 1).WillReturnRows(_rows)

	_mock.ExpectExec(`DELETE FROM "jwks" WHERE active = $1`).
		WithArgs(false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`UPDATE "jwks" SET "active"=$1 WHERE active = $2`).
		WithArgs(false, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateJWK(context.TODO(), _jwkOne)
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
			_, err := test.database.GetActiveJWK(context.TODO(), _jwkOne.Kid)
			if err != nil {
				t.Errorf("GetActiveJWK for %s returned err: %v", test.name, err)
			}

			_, err = test.database.GetActiveJWK(context.TODO(), _jwkTwo.Kid)
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

			_, err = test.database.GetActiveJWK(context.TODO(), _jwkOne.Kid)
			if err == nil {
				t.Errorf("GetActiveJWK for %s should have returned err", test.name)
			}

			_, err = test.database.GetActiveJWK(context.TODO(), _jwkTwo.Kid)
			if err == nil {
				t.Errorf("GetActiveJWK for %s should have returned err", test.name)
			}
		})
	}
}
